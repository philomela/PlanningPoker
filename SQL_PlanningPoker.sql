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
    NameRoom VARCHAR(200),
    Created DATETIME2,
    IsActive BIT,
    Creator INT,
    [GUID] VARCHAR(36)
)
GO 

/*Таблица пользователей*/
CREATE TABLE ServerPlanningPoker.Persons(
    Id INT PRIMARY KEY IDENTITY,
    LoginName VARCHAR(50) UNIQUE,
    Email VARCHAR(100) UNIQUE,
    [Password] VARCHAR(100),
    [Token] VARCHAR(36)   
)
GO

/*Таблица коннекшинов*/
CREATE TABLE ServerPlanningPoker.Connections(
    Id INT PRIMARY KEY,
    [GUID] VARCHAR(36),
    DateConnection DATETIME2,
    RoomId INT REFERENCES ServerPlanningPoker.Rooms (Id) ON DELETE CASCADE,
    PersonId INT REFERENCES ServerPlanningPoker.Persons (Id) ON DELETE CASCADE
)
GO 
 
/*Таблица наблюдаемого(Активные вью модели)*/
CREATE TABLE ServerPlanningPoker.ViewModels(
    Id INT PRIMARY KEY IDENTITY,
    Source VARCHAR(MAX),
    Created DATETIME2,
    RoomId INT REFERENCES ServerPlanningPoker.Rooms (Id) ON DELETE CASCADE
)
GO

/*Таблица задач*/
CREATE TABLE ServerPlanningPoker.Tasks(
    Id INT PRIMARY KEY IDENTITY,
    [Name] VARCHAR(MAX),
    Created DATETIME2,
    RoomId INT REFERENCES ServerPlanningPoker.Rooms (Id) ON DELETE CASCADE
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
                DECLARE @LastIdRoom INT;
                INSERT INTO ServerPlanningPoker.Rooms (NameRoom, Created, IsActive, Creator, [GUID])
                VALUES (@NameRoom, CURRENT_TIMESTAMP, 1, (SELECT Id FROM ServerPlanningPoker.Persons WHERE LoginName = @Creator),  NEWID());

                SELECT TOP(1) @LastIdRoom = Id FROM ServerPlanningPoker.Rooms ORDER BY Id DESC;

                        INSERT INTO ServerPlanningPoker.Tasks ([Name], Created, RoomId)
                        SELECT C.value('@name', 'nvarchar(max)'),
                               CURRENT_TIMESTAMP,
                               @LastIdRoom
                         FROM @Tasks.nodes('/tasks/task') T(C);               
                COMMIT;
                SELECT TOP(1) Rooms.[GUID] FROM ServerPlanningPoker.Rooms AS Rooms ORDER BY Id DESC;
            END TRY
            BEGIN CATCH

                ROLLBACK;
            END CATCH      
GO

/*Процедура создания и привязки коннекшенов к комнате*/
CREATE PROCEDURE [CreateConnection](@UUID VARCHAR(36), 
                                    @RoomId INT)
    AS
        BEGIN TRANSACTION
            BEGIN TRY
                INSERT INTO ServerPlanningPoker.Connections ([GUID], RoomId)
                VALUES (@UUID, @RoomId)
                COMMIT
            END TRY
            BEGIN CATCH
                ROLLBACK
            END CATCH
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


INSERT INTO ServerPlanningPoker.Persons(LoginName, Email, [Password], Token) 
VALUES('Roma2', 'romaphilomela@yandex2.ru', 'Rom@nkhik', NEWID())
                
EXECUTE CheckUser 'romaphilomela@yandex.ru', 'Rom@nkhik'
SELECT * FROM ServerPlanningPoker.Persons

EXECUTE Add_User 'philomeela', 'romaphilomela@yandex.ru', 'Rom@nkhik'

SELECT * FROM ServerPlanningPoker.Persons

DECLARE @ret VARCHAR(300);
EXEC [NewPlanningPokerRoom] 'RoomTdsdsdest', '
<tasks>
    <task name="Task1"></task>
    <task name="Task2"></task>
    <task name="Task3"></task>
</tasks>',
@ret OUTPUT
SELECT @ret


DROP PROCEDURE [NewPlanningPokerRoom]

DECLARE @temp XML = '
<tasks>
    <task name="TaskGoOne"></task>
    <task name="TaskGoTwo"></task>
    <task name="TaskGoThree"></task>
</tasks>'
SELECT C.value('@name', 'nvarchar(max)') FROM @temp.nodes('/tasks/task') T (C)


SELECT * FROM ServerPlanningPoker.Tasks ORDER BY 1 DESC
SELECT * FROM ServerPlanningPoker.Rooms ORDER BY 1 DESC
SELECT * FROM ServerPlanningPoker.ErrorsLog_Server 
SELECT * FROM ServerPlanningPoker.Connections
SELECT * FROM ServerPlanningPoker.ViewModels



SELECT 1 WHERE EXISTS ( SELECT [GUID] FROM FROM ServerPlanningPoker.Rooms WHERE GUID = '9c570ff5-bbec-40d8-b565-d4c373550dea')

DROP TABLE ServerPlanningPoker.Tasks
DROP TABLE ServerPlanningPoker.Rooms
DROP TABLE ServerPlanningPoker.ErrorsLog_Server 
DROP TABLE ServerPlanningPoker.Connections
DROP TABLE ServerPlanningPoker.ViewModels
DROP TABLE ServerPlanningPoker.Persons
DROP PROCEDURE u0932131_admin.NewPlanningPokerRoom
DROP PROCEDURE u0932131_admin.Add_User
DROP TABLE ServerPlanningPoker.Sessions


SELECT RegExp() FROM ServerPlanningPoker.Rooms 