/*Подключаемся к бд*/
USE u0932131_planningPoker
GO

/*Создаем схему в бд для системных таблиц сервера*/
CREATE SCHEMA ServerPlanningPoker AUTHORIZATION u0932131_admin
GO

/*Таблица ошибок сервера*/
CREATE TABLE ServerPlanningPoker.ErrorsLog_Server (
    Id INT PRIMARY KEY IDENTITY,
    ErrorText VARCHAR(MAX) NOT NULL,
    DateTimeError DATETIME2 NOT NULL,
)
GO

/*Хранимая процедура добавления ошибок в таблицу*/
CREATE PROCEDURE [SaveError](@ErrorText VARCHAR(MAX))
    AS
        BEGIN TRANSACTION
            BEGIN TRY
                INSERT INTO ServerPlanningPoker.ErrorsLog_Server(ErrorText, DateTimeError) 
                VALUES (@ErrorText, GETDATE())
                COMMIT;
            END TRY
            BEGIN CATCH
                ROLLBACK;
            END CATCH
GO

/*Создаем схему в бд для таблиц бизнес модели сервера*/
CREATE SCHEMA PlanningPokerBuisness AUTHORIZATION u0932131_admin
GO

/*Таблица комнат*/
CREATE TABLE ServerPlanningPoker.Rooms(
    Id INT PRIMARY KEY IDENTITY,
    NameRoom VARCHAR(200) NOT NULL,
    Created DATETIME2,
    IsActive BIT,
    Creator INT NOT NULL,
    [GUID] VARCHAR(36) UNIQUE
)
GO

/*Таблица пользователей*/
CREATE TABLE ServerPlanningPoker.Persons(
    Id INT PRIMARY KEY IDENTITY,
    LoginName VARCHAR(50) UNIQUE NOT NULL,
    Email VARCHAR(100) UNIQUE NOT NULL,
    [Password] VARCHAR(100),
    [Token] VARCHAR(36)   
)
GO

/*Таблица коннекшинов*/
CREATE TABLE ServerPlanningPoker.Connections(
    Id INT PRIMARY KEY IDENTITY,
    [GUID] VARCHAR(36),
    DateConnection DATETIME2,
    RoomId INT REFERENCES ServerPlanningPoker.Rooms (Id) ON DELETE CASCADE,
    PersonId INT REFERENCES ServerPlanningPoker.Persons (Id) ON DELETE CASCADE
)
GO 
 
/*Таблица наблюдаемого(Активные вью модели)*/
CREATE TABLE ServerPlanningPoker.ViewModels(
    Id INT PRIMARY KEY IDENTITY,
    Source XML,
    Created DATETIME2,
    RoomId INT REFERENCES ServerPlanningPoker.Rooms (Id) ON DELETE CASCADE
)
GO

/*Таблица задач*/
CREATE TABLE ServerPlanningPoker.Tasks(
    Id INT PRIMARY KEY IDENTITY,
    [Name] VARCHAR(MAX) NOT NULL,
    TimeDiscussion TINYINT NOT NULL CHECK (TimeDiscussion > 0 AND TimeDiscussion <= 10),
    Created DATETIME2,
    OnActive BIT NOT NULL,
    Completed BIT NOT NULL,
    RoomId INT REFERENCES ServerPlanningPoker.Rooms (Id) ON DELETE CASCADE
)
GO

/*Таблица голосов и оценок задач*/
CREATE TABLE ServerPlanningPoker.Votes(
    Id INT PRIMARY KEY IDENTITY,
    Vote BIT NOT NULL,
    Score INT NOT NULL,
    TaskId INT NOT NULL REFERENCES ServerPlanningPoker.Tasks (Id),
    RoomId INT NOT NULL REFERENCES ServerPlanningPoker.Rooms (Id) ON DELETE CASCADE,
    PersonId INT NOT NULL REFERENCES ServerPlanningPoker.Persons (Id)
)
GO

CREATE TABLE ServerPlanningPoker.TasksResults(
    Id INT PRIMARY KEY IDENTITY,
    Median DECIMAL(9,2) NOT NULL,
    DateCreated DATETIME2 NOT NULL,
    TaskId INT UNIQUE NOT NULL REFERENCES ServerPlanningPoker.Tasks (Id),
    RoomId INT NOT NULL REFERENCES ServerPlanningPoker.Rooms (Id) ON DELETE CASCADE
)
GO

/*Процедура добавление нового пользователя*/
CREATE PROCEDURE [Add_User](@LoginName VARCHAR(50), 
                            @Email VARCHAR(100), 
                            @Password VARCHAR(30))
    AS 
        BEGIN TRANSACTION        
                BEGIN TRY
                    IF EXISTS (SELECT Id FROM ServerPlanningPoker.Persons WHERE Email = @Email AND LoginName = @LoginName)
                    BEGIN
                        SELECT 'User with a pair of login/email exist'
                    END
                    ELSE IF EXISTS (SELECT Id FROM ServerPlanningPoker.Persons WHERE Email = @Email)
                    BEGIN
                        SELECT 'User with this email exists'
                    END
                    ELSE IF EXISTS (SELECT Id FROM ServerPlanningPoker.Persons WHERE LoginName = @LoginName)
                    BEGIN
                        SELECT 'User with this login exists'
                    END    
                    ELSE 
                    BEGIN
                    INSERT INTO ServerPlanningPoker.Persons(LoginName, Email, [Password], Token)
                        VALUES (@LoginName, @Email, @Password, NEWID());
                        SELECT 'Succsess'  
                    END           
                      COMMIT;              
                END TRY
                BEGIN CATCH
                    ROLLBACK;
                END CATCH
GO

/*Процедура формирования новой комнаты*/
CREATE PROCEDURE [NewPlanningPokerRoom](@NameRoom VARCHAR(200), 
                                        @Tasks XML,
                                        @Creator VARCHAR(50))
    AS  
        BEGIN TRANSACTION
            BEGIN TRY
                DECLARE @LastIdRoom INT, @RoomGUID VARCHAR(36) = NEWID();
                INSERT INTO ServerPlanningPoker.Rooms (NameRoom, Created, IsActive, Creator, [GUID])
                VALUES (@NameRoom, CURRENT_TIMESTAMP, 1, (SELECT Id FROM ServerPlanningPoker.Persons WHERE Email = @Creator),  @RoomGUID);

                SELECT TOP(1) @LastIdRoom = Id FROM ServerPlanningPoker.Rooms ORDER BY Id DESC;
 
                        INSERT INTO ServerPlanningPoker.Tasks ([Name], [TimeDiscussion], Created, OnActive, Completed, RoomId)
                        SELECT C.value('@name', 'nvarchar(max)'),
                               C.value('@time-discussion', 'tinyint'),
                               CURRENT_TIMESTAMP,
                               0,
                               0,
                               @LastIdRoom
                         FROM @Tasks.nodes('/tasks/task') T(C);                 
                COMMIT;
                SELECT TOP(1) Rooms.[GUID] FROM ServerPlanningPoker.Rooms AS Rooms ORDER BY Id DESC;
            END TRY
            BEGIN CATCH
                SELECT ERROR_MESSAGE()
                ROLLBACK;
            END CATCH
GO 
--SELECT * FROM ServerPlanningPoker.Tasks
/*Процедура создания и привязки коннекшенов к комнате*/
CREATE PROCEDURE [CreateConnection](@UUID VARCHAR(36), 
                                    @RoomGUID VARCHAR(36),
                                    @Email VARCHAR(100))
    AS
        BEGIN TRANSACTION
            BEGIN TRY
                DECLARE @TasksId TABLE (Id INT);
                DECLARE @PersonId INT = (SELECT Id FROM ServerPlanningPoker.Persons WHERE Email = @Email),
                        @RoomId INT = (SELECT Id FROM ServerPlanningPoker.Rooms WHERE GUID = @RoomGUID),
                        @xmlVM XML;
                        
                INSERT INTO @TasksId SELECT Id FROM ServerPlanningPoker.Tasks 
                    --INNER JOIN ServerPlanningPoker. SELECT * FROM ServerPlanningPoker.Tasks
                    WHERE RoomId = @RoomId --AND @PersonId = @PersonId;
                
                
                INSERT INTO ServerPlanningPoker.Connections ([GUID], DateConnection, RoomId, PersonId)
                VALUES (@UUID, 
                        CURRENT_TIMESTAMP,
                        (@RoomId),
                        (@PersonId)
                        );

                DECLARE @lastVal INT = (SELECT TOP(1) Id FROM @TasksId ORDER BY 1 DESC),
                        @currVal INT = (SELECT TOP(1) Id FROM @TasksId ORDER BY 1 ASC);
                WHILE (@currVal <= @lastVal)
                    BEGIN
                        IF NOT EXISTS (SELECT Id FROM ServerPlanningPoker.Votes 
                                                    WHERE PersonId = @PersonId AND RoomId = @RoomId AND TaskId = (SELECT @currVal))
                    BEGIN
                          INSERT INTO ServerPlanningPoker.Votes (Vote, Score, TaskId, RoomId, PersonId)
                          VALUES(0, 0, @currVal, @RoomId, @PersonId);
                    END
                    SET @currVal = @currVal + 1;
                END
                
                IF NOT EXISTS (SELECT Id FROM ServerPlanningPoker.ViewModels WHERE RoomId = @RoomId) --ПЕРЕСМОТРЕТЬ УСЛОВИЕ, ВОЗМОЖНО ОБНОВЛЯТЬ ВЬЮ МОДЕЛЬ
                    BEGIN 
                          EXECUTE [Build_First_ViewModel] @RoomId, @xmlVM OUTPUT;
                          
                          INSERT INTO ServerPlanningPoker.ViewModels (Source, Created, RoomId)
                          VALUES (
                                    @xmlVM,
                                    CURRENT_TIMESTAMP,
                                    @RoomId
                                ); 
                    END
                ELSE 
                    BEGIN
                        EXECUTE [Build_First_ViewModel] @RoomId, @xmlVM OUTPUT;
                        
                        UPDATE ServerPlanningPoker.ViewModels
                        SET Source = @xmlVM
                        WHERE RoomId = @RoomId;
                    END
                COMMIT;
            END TRY
            BEGIN CATCH
                ROLLBACK
            END CATCH
GO 
/*Процедура получения View Model*/
CREATE PROCEDURE [Build_First_ViewModel] (@roomId INT, @xmlVMOut XML OUTPUT)
    AS
        SET @xmlVMOut = (SELECT (SELECT P.LoginName AS '@UserName',
                                        (ROW_NUMBER() OVER (ORDER BY P.Id ASC)) AS '@Id' 
                                FROM ServerPlanningPoker.Persons AS P
                                WHERE P.Id IN (SELECT C.PersonId FROM ServerPlanningPoker.Connections AS C WHERE C.RoomId = @roomId)
                                FOR XML PATH ('Person'), TYPE ) AS Persons,
                                (SELECT (ROW_NUMBER() OVER (ORDER BY T.Id ASC)) AS '@Id',
                                        T.Name AS '@NameTask',
                                        T.OnActive AS '@IsCurrentActive',
                                        T.TimeDiscussion AS '@TimeDiscussion',
                                        (SELECT Median
                                        FROM ServerPlanningPoker.TasksResults
                                        WHERE RoomId = @RoomId AND TaskId = T.Id) AS '@Median',
                                        (SELECT ROW_NUMBER() OVER (ORDER BY V.PersonId ASC) AS '@PersonId',
                                                V.Vote AS '@Vote',
                                                V.Score AS '@Score'
                                        FROM ServerPlanningPoker.Votes AS V 
                                        WHERE RoomId = @roomId AND TaskId = T.Id
                                        FOR XML PATH ('PersonTask'), TYPE)
                                FROM ServerPlanningPoker.Tasks AS T 
                                WHERE T.RoomId = @roomId
                                FOR XML PATH ('Task'), TYPE) AS Tasks,
                                R.NameRoom
                        FROM ServerPlanningPoker.Rooms AS R
                     WHERE R.Id = @roomId
                     FOR XML RAW ('Room'), TYPE);
GO

/*Процедура проверки пользователя*/
CREATE PROCEDURE [CheckUser] (@email VARCHAR(100), 
                              @password VARCHAR(100))
    AS  
        DECLARE @ResultCheck bit;
            IF EXISTS (SELECT id FROM ServerPlanningPoker.Persons WHERE Email = @email AND [Password] = @password)
                BEGIN 
                    SET @ResultCheck = 1; 
                END
            ELSE
                BEGIN
                    SET @ResultCheck = 0;
                END
        SELECT @ResultCheck;
GO

/*Проверка создателя комнаты*/
CREATE PROCEDURE [CheckCreator] (@email VARCHAR(100),
                                @roomUID VARCHAR(36))
    AS
        IF (@email = (SELECT Email FROM ServerPlanningPoker.Persons 
                                        WHERE Id = (SELECT Creator FROM ServerPlanningPoker.Rooms 
                                                            WHERE GUID = @roomUID)))
            BEGIN 
                SELECT 'Creator';
            END
        ELSE 
            BEGIN 
                SELECT 'User'
            END
GO

--EXEC Push_And_Get_Changes  '<Change><AddVote vote="1" /></Change>', 'add_vote', '58953a61-abb7-425e-9441-1ca92b97b6d1', 'romaphilomela@yandex.ru'
/*Хранимая процедура сохраниения и получения xml ViewModel*/ --select * from ServerPlanningPoker.votes order by 1 desc
CREATE PROCEDURE [Push_And_Get_Changes] (@xmlChanges XML, 
                                        @nameChanges VARCHAR(50), 
                                        @roomGUID VARCHAR(36), 
                                        @email VARCHAR(50))
    AS
        BEGIN TRANSACTION 
        BEGIN TRY
            
            DECLARE @xmlOut XML,
                    @RoomId INT = (SELECT Id FROM ServerPlanningPoker.Rooms WHERE GUID = @roomGUID),
                    @PersonId INT = (SELECT Id FROM ServerPlanningPoker.Persons WHERE Email = @email),
                    @TaskId INT;
                
                IF @nameChanges = 'ChangeVote' 
                    BEGIN
                        DECLARE @Score INT, @Vote INT;
                    
                        SELECT @Score = C.value('@score', 'INT')
                                ,@Vote = C.value('@vote', 'INT')                     
                            FROM @xmlChanges.nodes('/Change/AddVote') T(C);

                        SET @TaskId = (SELECT Id 
                                            FROM ServerPlanningPoker.Tasks 
                                                WHERE RoomId = @RoomId AND OnActive = 1)

                        UPDATE ServerPlanningPoker.Votes
                        SET Score = @Score, Vote = 1
                        WHERE RoomId = @RoomId AND PersonId = @PersonId AND TaskId = @TaskId;
                        
                        EXEC Build_First_ViewModel @RoomId, @xmlOut OUTPUT;
                        SELECT @xmlOut;
                    END
                
                ELSE IF @nameChanges = 'ChangeGetVM'
                    BEGIN
                        EXEC Build_First_ViewModel @RoomId, @xmlOut OUTPUT; 
                        SELECT @xmlOut;
                    END

                ELSE IF @nameChanges = 'StartVoting'
                    BEGIN

                        WITH 
                        tasksTable AS (
                            SELECT ROW_NUMBER() OVER (PARTITION BY RoomId ORDER BY Id ASC) AS RowId
                                   ,Id
                                FROM ServerPlanningPoker.Tasks WHERE RoomId = @RoomId
                        )
                        UPDATE ServerPlanningPoker.Tasks
                        SET OnActive = 1 
                            WHERE RoomId = @RoomId AND Id = (SELECT Id 
                                                                FROM tasksTable 
                                                                    WHERE RowId = (SELECT TOP(1) C.value('@taskId', 'INT')
                                                                                        FROM @xmlChanges.nodes('/Change/StartVoting') T(C)));

                        EXEC Build_First_ViewModel @RoomId, @xmlOut OUTPUT; 
                        SELECT @xmlOut;
                    END 
                ELSE IF @nameChanges = 'StopVoting'
                    BEGIN
                        
                        UPDATE ServerPlanningPoker.Tasks
                            SET OnActive = 0, Completed = 1
                                WHERE RoomId = @RoomId AND OnActive = 1;
                        
                        INSERT INTO ServerPlanningPoker.TasksResults(Median, DateCreated, TaskId, RoomId)
                        VALUES ((SELECT TOP(1) PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY Score)
                                                            OVER (PARTITION BY RoomId) AS Median
                                    FROM ServerPlanningPoker.Votes 
                                        WHERE TaskId = (SELECT TOP(1) Id 
                                                            FROM ServerPlanningPoker.Tasks 
                                                                WHERE RoomId = @RoomId AND Completed = 1 AND OnActive = 0 ORDER BY Id DESC))
                                ,CURRENT_TIMESTAMP
                                ,(SELECT TOP(1) Id 
                                    FROM ServerPlanningPoker.Tasks
                                        WHERE RoomId = @RoomId AND Completed = 1 AND OnActive = 0 ORDER BY Id DESC)
                                ,@RoomId
                                )

                        EXEC Build_First_ViewModel @RoomId, @xmlOut OUTPUT; 
                        SELECT @xmlOut;
                    END
                COMMIT;
            END TRY

            BEGIN CATCH
                
                ROLLBACK;
            END CATCH
GO

SELECT * FROM ServerPlanningPoker.Votes ORDER BY RoomId DESC
SELECT * FROM ServerPlanningPoker.TasksResults ORDER BY 1 DESC
DELETE FROM ServerPlanningPoker.TasksResults

SELECT 
INSERT INTO ServerPlanningPoker.Persons(LoginName, Email, [Password], Token) 
VALUES('Roma2', 'romaphilomela@yandex2.ru', 'Rom@nkhik', NEWID())
                
EXECUTE CheckUser 'romaphilomela@yandex.ru', 'Rom@nkhik'
SELECT * FROM ServerPlanningPoker.Persons

EXECUTE Add_User 'philomelka', 'philomela@yandex.ru', 'Rom@nkhik'

SELECT * FROM ServerPlanningPoker.Persons

DECLARE @ret VARCHAR(300);
EXEC [NewPlanningPokerRoom] 'RoomTdsdsdest', '
<tasks>
    <task name="Task1" time-discussion="5"></task>
    <task name="Task2" time-discussion="5"></task>
    <task name="Task3" time-discussion="6"></task>
</tasks>', 'romaphilomela@yandex.ru'
@ret OUTPUT
SELECT @ret 


DROP PROCEDURE [NewPlanningPokerRoom]

DECLARE @temp XML = '
<tasks>
    <task name="TaskGoOne" time-discussion="5"></task>
    <task name="TaskGoTwo" time-discussion="6"></task>
    <task name="TaskGoThree" time-discussion="7"></task>
</tasks>'
SELECT C.value('@name', 'nvarchar(max)'), C.value('@time-discussion', 'tinyint') FROM @temp.nodes('/tasks/task') T (C)

INSERT INTO ServerPlanningPoker.Rooms (NameRoom, Created, IsActive, Creator, [GUID]) VALUES ('room01', '2020-12-10', 1, 0, 'huid')
INSERT INTO ServerPlanningPoker.Tasks ([Name], [TimeDiscussion], Created, OnActive, RoomId) VALUES ('Hello', 5, '2020-12-10', 0, 6)

SELECT * FROM ServerPlanningPoker.Tasks ORDER BY 1 DESC
SELECT * FROM ServerPlanningPoker.Rooms ORDER BY 1 DESC
SELECT * FROM ServerPlanningPoker.ErrorsLog_Server 
SELECT * FROM ServerPlanningPoker.Connections
SELECT * FROM ServerPlanningPoker.ViewModels
SELECT * FROM ServerPlanningPoker.Persons



SELECT 1 WHERE EXISTS ( SELECT [GUID] FROM FROM ServerPlanningPoker.Rooms WHERE GUID = '9c570ff5-bbec-40d8-b565-d4c373550dea')

DROP TABLE ServerPlanningPoker.Votes
DROP TABLE ServerPlanningPoker.Tasks
DROP TABLE ServerPlanningPoker.Rooms
DROP TABLE ServerPlanningPoker.ErrorsLog_Server 
DROP TABLE ServerPlanningPoker.Connections
DROP TABLE ServerPlanningPoker.ViewModels
DROP TABLE ServerPlanningPoker.Persons
DROP TABLE ServerPlanningPoker.Sessions
DROP TABLE ServerPlanningPoker.TasksResults
DROP PROCEDURE u0932131_admin.NewPlanningPokerRoom
DROP PROCEDURE u0932131_admin.Add_User
DROP PROCEDURE u0932131_admin.CreateConnection
DROP PROCEDURE Get_ViewModel
DROP PROCEDURE Push_And_Get_Changes
DROP PROCEDURE Build_First_ViewModel


SELECT RegExp() FROM ServerPlanningPoker.Rooms 
SELECT * FROM ServerPlanningPoker.Persons ORDER BY 1 DESC




CREATE TABLE TestTable_01 (
    id INT PRIMARY KEY IDENTITY,
    surname VARCHAR(100) NULL,
    [name] VARCHAR(100) NULL,
    number_flight CHAR(10) NULL
)

INSERT INTO TestTable_01 ([name], surname, number_flight) VALUES ('Soloviev', 'Roman', 'A')

INSERT INTO TestTable_01 VALUE ()

SELECT * FROM TestTable_01

DELETE FROM ServerPlanningPoker.Persons
