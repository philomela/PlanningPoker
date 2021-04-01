/*Connect DB*/
USE u0932131_planningPoker
GO

/*Create schema in db for system tables*/
CREATE SCHEMA ServerPlanningPoker AUTHORIZATION u0932131_admin
GO

/*Errors table*/
CREATE TABLE ServerPlanningPoker.ErrorsLog_Server (
    Id INT PRIMARY KEY IDENTITY,
    ErrorText VARCHAR(MAX) NOT NULL,
    DateTimeError DATETIME2 NOT NULL,
)
GO

/*Proc for add errors in error table*/
CREATE PROCEDURE ServerPlanningPoker.[SaveError](@ErrorText VARCHAR(MAX))
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

/*Rooms table*/
CREATE TABLE ServerPlanningPoker.Rooms(
    Id INT PRIMARY KEY IDENTITY,
    NameRoom VARCHAR(200) NOT NULL,
    Created DATETIME2,
    Deleted DATETIME2,
    IsActive BIT,
    Creator INT NOT NULL,
    [GUID] VARCHAR(36) UNIQUE
)
GO

/*Users table*/
CREATE TABLE ServerPlanningPoker.Persons(
    Id INT PRIMARY KEY IDENTITY,
    LoginName VARCHAR(50) UNIQUE NOT NULL,
    Email VARCHAR(100) UNIQUE NOT NULL,
    [Password] VARCHAR(100),
    [Token] VARCHAR(36)   
)
GO

/*Connections table*/
CREATE TABLE ServerPlanningPoker.Connections(
    Id INT PRIMARY KEY IDENTITY,
    [GUID] VARCHAR(36),
    DateConnection DATETIME2,
    RoomId INT REFERENCES ServerPlanningPoker.Rooms (Id) ON DELETE CASCADE,
    PersonId INT REFERENCES ServerPlanningPoker.Persons (Id) ON DELETE CASCADE
)
GO 
 
/*ViewModels table*/
CREATE TABLE ServerPlanningPoker.ViewModels(
    Id INT PRIMARY KEY IDENTITY,
    Source XML,
    Created DATETIME2,
    RoomId INT REFERENCES ServerPlanningPoker.Rooms (Id) ON DELETE CASCADE
)
GO

/*Tasks table*/
CREATE TABLE ServerPlanningPoker.Tasks(
    Id INT PRIMARY KEY IDENTITY,
    [Name] VARCHAR(MAX) NOT NULL,
    TimeDiscussion TINYINT NOT NULL CHECK (TimeDiscussion > 0 AND TimeDiscussion <= 10),
    Created DATETIME2,
    OnActive BIT NOT NULL,
    Completed BIT NOT NULL,
    DateComplited DATETIME2 NOT NULL DEFAULT CURRENT_TIMESTAMP,
    RoomId INT REFERENCES ServerPlanningPoker.Rooms (Id) ON DELETE CASCADE
)
GO

/*Votes and voices table*/
CREATE TABLE ServerPlanningPoker.Votes(
    Id INT PRIMARY KEY IDENTITY,
    Vote BIT NOT NULL,
    Score INT NOT NULL,
    TaskId INT NOT NULL REFERENCES ServerPlanningPoker.Tasks (Id),
    RoomId INT NOT NULL REFERENCES ServerPlanningPoker.Rooms (Id) ON DELETE CASCADE,
    PersonId INT NOT NULL REFERENCES ServerPlanningPoker.Persons (Id)
)
GO

/*Task results table*/
CREATE TABLE ServerPlanningPoker.TasksResults(
    Id INT PRIMARY KEY IDENTITY,
    Median DECIMAL(9,2),
    DateCreated DATETIME2 NOT NULL,
    TaskId INT UNIQUE NOT NULL REFERENCES ServerPlanningPoker.Tasks (Id),
    RoomId INT NOT NULL REFERENCES ServerPlanningPoker.Rooms (Id) ON DELETE CASCADE
)
GO

/*Restored accounts table*/
CREATE TABLE ServerPlanningPoker.RestoredAccounts(
    Id INT PRIMARY KEY IDENTITY,
    DateRequest DATETIME2 NOT NULL,
    Link VARCHAR(100) NOT NULL,
    IsActive BIT NOT NULL,
    AccountId INT NOT NULL FOREIGN KEY REFERENCES ServerPlanningPoker.Persons (Id) ON DELETE CASCADE
)
GO

/*Score-dictionary table*/
CREATE TABLE ServerPlanningPoker.Scores_Dictionary (
    Id INT PRIMARY KEY IDENTITY,
    Score INT NOT NULL,
    [Decryption] VARCHAR(300)
)
GO

/*Init insert score-dictionary table*/
INSERT INTO ServerPlanningPoker.Scores_Dictionary(Score, [Decryption])
        VALUES(999, 'Coffee break')
              ,(777, 'Question')
              ,(0, 'Not voice')
GO

/*Proc add new user*/
CREATE PROCEDURE ServerPlanningPoker.[Add_User](@LoginName VARCHAR(50), 
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

/*Proc restore account*/
CREATE PROCEDURE ServerPlanningPoker.[RestoreAccount] (@Link VARCHAR(100),
                                   @Password VARCHAR(100))
    AS
        BEGIN TRANSACTION
            BEGIN TRY 
                IF 1 = (SELECT IsActive FROM ServerPlanningPoker.RestoredAccounts WHERE Link = @Link) 
                    BEGIN             
                        UPDATE ServerPlanningPoker.Persons
                        SET [Password] = @Password
                        WHERE Id = (SELECT TOP(1) AccountId 
                                                    FROM ServerPlanningPoker.RestoredAccounts 
                                                        WHERE Link = @Link ORDER BY Id DESC);
                        UPDATE ServerPlanningPoker.RestoredAccounts
                        SET IsActive = 0
                        WHERE Link = @Link;
                    SELECT 1;
                    END
                ELSE SELECT 0;
                COMMIT;
            END TRY 

            BEGIN CATCH
                SELECT 0;
                ROLLBACK;
            END CATCH
GO

/*Proc create link for restore account*/
CREATE PROCEDURE ServerPlanningPoker.[CreateAccountRecoveryLink] (@Email VARCHAR(100))
    AS
        BEGIN TRANSACTION
            BEGIN TRY     
                DECLARE @NewLink VARCHAR(36) = NEWID();

                MERGE ServerPlanningPoker.RestoredAccounts AS RS_BASE
                USING (SELECT CURRENT_TIMESTAMP AS DateRequest, 
                              NEWID() AS Link, 
                              Id AS AccountId 
                                    FROM ServerPlanningPoker.Persons 
                                        WHERE Email = @Email) AS RS_SOURCE
                ON (RS_SOURCE.AccountId = RS_BASE.AccountId)
                WHEN MATCHED  THEN 
                    UPDATE SET RS_BASE.Link = @NewLink,
                               RS_BASE.IsActive = 1,
                               RS_BASE.DateRequest = CURRENT_TIMESTAMP
                WHEN NOT MATCHED THEN 
                    INSERT (DateRequest, Link, IsActive, AccountId)
                    VALUES (CURRENT_TIMESTAMP, @NewLink, 1, RS_SOURCE.AccountId); 
                
                SELECT Link 
                        FROM ServerPlanningPoker.RestoredAccounts
                            WHERE AccountId = (SELECT Id 
                                                        FROM ServerPlanningPoker.Persons 
                                                            WHERE Email = @Email);
                COMMIT;
            END TRY 

            BEGIN CATCH
                ROLLBACK;
            END CATCH
GO

/*Proc add new room*/
CREATE PROCEDURE ServerPlanningPoker.[NewPlanningPokerRoom](@NameRoom VARCHAR(200), 
                                        @Tasks XML,
                                        @Creator VARCHAR(50))
    AS  
        BEGIN TRANSACTION
            BEGIN TRY
                DECLARE @LastIdRoom INT, @RoomGUID VARCHAR(36) = NEWID();

                IF(@NameRoom = '' OR @NameRoom IS NULL)
                BEGIN
                    ROLLBACK
                END 

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
                ROLLBACK;
            END CATCH
GO 

/*Proc get View Model*/
CREATE PROCEDURE ServerPlanningPoker.[Get_ViewModel] (@roomId INT, @xmlVMOut XML OUTPUT)
    AS
        SET @xmlVMOut = ISNULL((SELECT (SELECT P.LoginName AS '@UserName',
                                        (ROW_NUMBER() OVER (ORDER BY P.Id ASC)) AS '@Id' 
                                FROM ServerPlanningPoker.Persons AS P
                                WHERE P.Id IN (SELECT C.PersonId FROM ServerPlanningPoker.Connections AS C WHERE C.RoomId = @roomId)
                                FOR XML PATH ('Person'), TYPE ) AS Persons,
                                (SELECT (ROW_NUMBER() OVER (ORDER BY T.Id ASC)) AS '@Id',
                                        T.Name AS '@NameTask',
                                        T.Completed AS '@Completed',
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
                     WHERE R.Id = @roomId AND R.IsActive = 1
                     FOR XML RAW ('Room'), TYPE), (SELECT 'UnknownRoom' AS 'Error' FOR XML RAW ('Room'), TYPE));
GO

/*Proc create connection and linked their with room*/
CREATE PROCEDURE ServerPlanningPoker.[CreateConnection](@UUID VARCHAR(36), 
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
                    WHERE RoomId = @RoomId
                
                
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
                
                IF NOT EXISTS (SELECT Id FROM ServerPlanningPoker.ViewModels WHERE RoomId = @RoomId)
                    BEGIN 
                          EXECUTE ServerPlanningPoker.[Get_ViewModel] @RoomId, @xmlVM OUTPUT;
                          
                          INSERT INTO ServerPlanningPoker.ViewModels (Source, Created, RoomId)
                          VALUES (
                                    @xmlVM,
                                    CURRENT_TIMESTAMP,
                                    @RoomId
                                ); 
                    END
                ELSE 
                    BEGIN
                        EXECUTE ServerPlanningPoker.[Get_ViewModel] @RoomId, @xmlVM OUTPUT;
                        
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

/*Proc check user*/
CREATE PROCEDURE ServerPlanningPoker.[CheckUser] (@email VARCHAR(100), 
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

/*Proc check creator room*/
CREATE PROCEDURE ServerPlanningPoker.[CheckCreator] (@email VARCHAR(100),
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

/*Proc save and get xml ViewModel*/
CREATE PROCEDURE ServerPlanningPoker.[Push_And_Get_Changes] (@xmlChanges XML, 
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
                        SET Score = @Score, Vote = @Vote
                        WHERE RoomId = @RoomId AND PersonId = @PersonId AND TaskId = @TaskId;
                        
                        EXEC ServerPlanningPoker.[Get_ViewModel] @RoomId, @xmlOut OUTPUT;
                        SELECT @xmlOut;
                    END
                
                ELSE IF @nameChanges = 'ChangeGetVM'
                    BEGIN
                        EXEC ServerPlanningPoker.[Get_ViewModel] @RoomId, @xmlOut OUTPUT; 
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
                        SET OnActive = 1,
                            DateComplited = CURRENT_TIMESTAMP
                            WHERE RoomId = @RoomId AND Id = (SELECT Id 
                                                                FROM tasksTable 
                                                                    WHERE RowId = (SELECT TOP(1) C.value('@taskId', 'INT')
                                                                                        FROM @xmlChanges.nodes('/Change/StartVoting') T(C)))

                        EXEC ServerPlanningPoker.[Get_ViewModel] @RoomId, @xmlOut OUTPUT; 
                        SELECT @xmlOut;
                    END 
                ELSE IF @nameChanges = 'StopVoting'
                    BEGIN
                        IF EXISTS(SELECT 1 FROM ServerPlanningPoker.Tasks WHERE RoomId = @RoomId AND OnActive = 1 AND Completed = 0)
                            BEGIN
                                UPDATE ServerPlanningPoker.Tasks
                                    SET OnActive = 0, Completed = 1 , DateComplited = CURRENT_TIMESTAMP
                                        WHERE Id = (SELECT TOP(1)Id 
                                                            FROM ServerPlanningPoker.Tasks 
                                                                WHERE RoomId = @RoomId 
                                                                AND OnActive = 1 ORDER BY DateComplited DESC);
                                
                                INSERT INTO ServerPlanningPoker.TasksResults(Median, DateCreated, TaskId, RoomId)
                                VALUES ((SELECT TOP(1) PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY Score)
                                                                    OVER (PARTITION BY RoomId) AS Median
                                            FROM ServerPlanningPoker.Votes 
                                                WHERE TaskId = (SELECT TOP(1) Id 
                                                                    FROM ServerPlanningPoker.Tasks ORDER BY DateComplited DESC)
                                                    AND Score NOT IN (SELECT Score FROM ServerPlanningPoker.Scores_Dictionary))
                                        ,CURRENT_TIMESTAMP
                                        ,(SELECT TOP(1) Id 
                                            FROM ServerPlanningPoker.Tasks ORDER BY DateComplited DESC)
                                        ,@RoomId
                                        ) 
                            END
                        EXEC ServerPlanningPoker.[Get_ViewModel] @RoomId, @xmlOut OUTPUT; 
                        SELECT @xmlOut;
                    END
                ELSE IF @nameChanges = 'FinishPlanning'
                    BEGIN
                        
                        UPDATE ServerPlanningPoker.Rooms 
                        SET Deleted = CURRENT_TIMESTAMP,
                            IsActive = 0
                        WHERE Id = @RoomId

                        EXEC ServerPlanningPoker.[Get_ViewModel] @RoomId, @xmlOut OUTPUT; 
                        SELECT @xmlOut;
                    END
                COMMIT;

            END TRY

            BEGIN CATCH
                ROLLBACK;
            END CATCH
GO

/*Proc check user email*/
CREATE PROCEDURE ServerPlanningPoker.[Check_User_Email] (@Email VARCHAR(50))
    AS 
        BEGIN TRANSACTION 
        BEGIN TRY 
            SELECT Email 
                    FROM ServerPlanningPoker.Persons WITH(nolock) 
                        WHERE Email = @email
            COMMIT;
        END TRY

        BEGIN CATCH
            ROLLBACK;
        END CATCH
GO



DROP TABLE ServerPlanningPoker.[Votes]
DROP TABLE ServerPlanningPoker.[Tasks]
DROP TABLE ServerPlanningPoker.[Rooms]
DROP TABLE ServerPlanningPoker.[ErrorsLog_Server] 
DROP TABLE ServerPlanningPoker.[Connections]
DROP TABLE ServerPlanningPoker.[ViewModels]
DROP TABLE ServerPlanningPoker.[Persons]
DROP TABLE ServerPlanningPoker.[TasksResults]
DROP TABLE ServerPlanningPoker.[Scores_Dictionary]
DROP TABLE ServerPlanningPoker.[RestoredAccounts]
DROP PROCEDURE ServerPlanningPoker.[NewPlanningPokerRoom]
DROP PROCEDURE ServerPlanningPoker.[Add_User]
DROP PROCEDURE ServerPlanningPoker.[CreateConnection]
DROP PROCEDURE ServerPlanningPoker.[Get_ViewModel]
DROP PROCEDURE ServerPlanningPoker.[Push_And_Get_Changes]
DROP PROCEDURE ServerPlanningPoker.[Get_Password]
DROP PROCEDURE ServerPlanningPoker.[CreateAccountRecoveryLink]
DROP PROCEDURE ServerPlanningPoker.[RestoreAccount]
DROP PROCEDURE ServerPlanningPoker.[Check_User_Email]
DROP PROCEDURE ServerPlanningPoker.[CheckCreator]
DROP PROCEDURE ServerPlanningPoker.[CheckUser]
DROP PROCEDURE ServerPlanningPoker.[SaveError]
DROP SCHEMA [ServerPlanningPoker]

SELECT * FROM ServerPlanningPoker.RestoredAccounts
SELECT * FROM ServerPlanningPoker.Persons
