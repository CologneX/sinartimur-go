-- Enable UUID Extension
Create Extension If Not Exists "uuid-ossp";

-- Table: Admin
Create Table Appuser
(
    Id            Uuid Primary Key Default Uuid_Generate_V4(),
    Username      VARCHAR(100) Unique Not Null,
    Password_Hash TEXT                Not Null,
    Is_Active     BOOLEAN          Default True,
    Is_Admin      BOOLEAN          Default False,
    Is_Hr         BOOLEAN          Default False,
    Is_Finance    BOOLEAN          Default False,
    Is_Inventory  BOOLEAN          Default False,
    Is_Sales      BOOLEAN          Default False,
    Is_Purchase   BOOLEAN          Default False,
    Created_At    Timestamptz      Default Current_Timestamp,
    Updated_At    Timestamptz      Default Current_Timestamp
);

-- Create Table Roles
-- (
--     Id          Uuid Primary Key Default Uuid_Generate_V4(),
--     Name        VARCHAR(50) Unique Not Null,
--     Description TEXT,
--     Created_At  Timestamptz      Default Current_Timestamp,
--     Updated_At  Timestamptz      Default Current_Timestamp
-- );
--
-- Create Table User_Roles
-- (
--     Id          Uuid Primary Key Default Uuid_Generate_V4(),
--     User_Id     Uuid References Users (Id) On Delete Cascade,
--     Role_Id     Uuid References Roles (Id) On Delete Cascade,
--     Assigned_At Timestamptz      Default Current_Timestamp,
--     Unique (User_Id, Role_Id)
-- );

-- Table: HR
Create Table Employee
(
    Id         Uuid Primary Key      Default Uuid_Generate_V4(),
    Name       VARCHAR(150) Not Null,
    Position   VARCHAR(100) Not Null,
    Phone      VARCHAR(20)  Not Null,
    Nik        VARCHAR(20)  Not Null,
    Hired_Date Timestamptz  Not Null Default Current_Timestamp,
    Created_At Timestamptz           Default Current_Timestamp,
    Updated_At Timestamptz           Default Current_Timestamp,
    Deleted_At Timestamptz           Default Null
);


Create Table Wage
(
    Id           Uuid Primary Key Default Uuid_Generate_V4(),
    Employee_Id  Uuid References Employee (Id) On Delete Cascade,
    Total_Amount NUMERIC(12, 2) Not Null,
    Month        INT            Not Null,
    Year         INT            Not Null,
    Created_At   Timestamptz      Default Current_Timestamp,
    Updated_At   Timestamptz      Default Current_Timestamp,
    Deleted_At   Timestamptz      Default Null
);

Create Table Wage_Detail
(
    Id             Uuid Primary Key Default Uuid_Generate_V4(),
    Wage_Id        Uuid References Wage (Id) On Delete Cascade,
    Component_Name VARCHAR(100)   Not Null,
    Description    TEXT,
    Amount         NUMERIC(12, 2) Not Null,
    Created_At     Timestamptz      Default Current_Timestamp,
    Updated_At     Timestamptz      Default Current_Timestamp,
    Deleted_At     Timestamptz      Default Null
);

-- Table: Financial Transactions
Create Table Financial_Transaction
(
    Id               Uuid Primary Key Default Uuid_Generate_V4(),
    User_Id          Uuid References Appuser (Id),
    Amount           NUMERIC(15, 2) Not Null,
    Type             VARCHAR(50)    Not Null,
    Description      TEXT,
    Transaction_Date Timestamptz    Not Null,
    Created_At       Timestamptz      Default Current_Timestamp,
    Edited_At        Timestamptz      Default Null,
    Deleted_At       Timestamptz      Default Null
);

-- Table: Inventory
Create Table Category
(
    Id          Uuid Primary Key Default Uuid_Generate_V4(),
    Name        VARCHAR(255) Not Null,
    Description TEXT,
    Created_At  Timestamptz      Default Current_Timestamp,
    Updated_At  Timestamptz      Default Current_Timestamp,
    Deleted_At  Timestamptz      Default Null
);

Create Table Unit
(
    Id          Uuid Primary Key Default Uuid_Generate_V4(),
    Name        VARCHAR(50) Not Null,
    Description TEXT,
    Created_At  Timestamptz      Default Current_Timestamp,
    Updated_At  Timestamptz      Default Current_Timestamp,
    Deleted_At  Timestamptz      Default Null
);

Create Table Product
(
    Id          Uuid Primary Key Default Uuid_Generate_V4(),
    Name        VARCHAR(255)   Not Null,
    Description TEXT,
    Price       NUMERIC(15, 2) Not Null,
    Category_Id Uuid             Default Null References Category (Id) On Delete Cascade,
    Unit_Id     Uuid             Default Null References Unit (Id) On Delete Cascade,
    Created_At  Timestamptz      Default Current_Timestamp,
    Updated_At  Timestamptz      Default Current_Timestamp,
    Deleted_At  Timestamptz      Default Null
);


Create Table Storage
(
    Id         Uuid Primary Key Default Uuid_Generate_V4(),
    Name       VARCHAR(255) Not Null,
    Location   TEXT         Not Null,
    Created_At Timestamptz      Default Current_Timestamp,
    Updated_At Timestamptz      Default Current_Timestamp,
    Deleted_At Timestamptz      Default Null
);

Create Table Inventory
(
    Id               Uuid Primary Key Default Uuid_Generate_V4(),
    Product_Id       Uuid References Product (Id) On Delete Cascade,
    Storage_Id       Uuid References Storage (Id) On Delete Cascade,
    Quantity         INT Not Null     Default 0,
    Minimum_Quantity INT Not Null     Default 0,
    Created_At       Timestamptz      Default Current_Timestamp,
    Updated_At       Timestamptz      Default Current_Timestamp
);


Create Table Inventory_Log
(
    Id           Uuid Primary Key Default Uuid_Generate_V4(),
    Inventory_Id Uuid References Inventory (Id) On Delete Cascade,
    User_Id      Uuid References Appuser (Id),
    Action       VARCHAR(50) Not Null, -- e.g., "add", "remove", "transfer"
    Quantity     INT         Not Null,
    Log_Date     Timestamptz      Default Current_Timestamp,
    Description  TEXT,
    Created_At   Timestamptz      Default Current_Timestamp,
    Updated_At   Timestamptz      Default Current_Timestamp
);

-- Table: Purchase
Create Table Supplier
(
    Id         Uuid Primary Key Default Uuid_Generate_V4(),
    Name       VARCHAR(255) Not Null,
    Address    TEXT,
    Telephone  VARCHAR(50),
    Created_At Timestamptz      Default Current_Timestamp,
    Updated_At Timestamptz      Default Current_Timestamp,
    Deleted_At Timestamptz      Default Null
);

Create Table Purchase_Order
(
    Id           Uuid Primary Key Default Uuid_Generate_V4(),
    Supplier_Id  Uuid             Default Null References Supplier (Id) On Delete Set Null,
    Order_Date   Timestamptz      Default Current_Timestamp,
    Status       VARCHAR(50)                   Not Null,
    Total_Amount NUMERIC(15, 2)                Not Null,
    Created_By   Uuid                          Not Null References Appuser(Id) On Delete Set Null,
    Created_At   Timestamptz      Default Current_Timestamp,
    Updated_At   Timestamptz      Default Current_Timestamp,
    Cancelled_At Timestamptz      Default Null
);

Create Table Purchase_Order_Detail
(
    Id                Uuid Primary Key Default Uuid_Generate_V4(),
    Purchase_Order_Id Uuid References Purchase_Order (Id) On Delete Cascade,
    Product_Id        Uuid References Product (Id) On Delete Cascade,
    Quantity          INT            Not Null,
    Price             NUMERIC(15, 2) Not Null,
    Created_At        Timestamptz      Default Current_Timestamp,
    Updated_At        Timestamptz      Default Current_Timestamp
);

-- Table: Sales
Create Table Sales_Order
(
    Id            Uuid Primary Key Default Uuid_Generate_V4(),
    Customer_Name VARCHAR(255)   Not Null,
    Order_Date    Timestamptz      Default Current_Timestamp,
    Status        VARCHAR(50)    Not Null,
    Total_Amount  NUMERIC(15, 2) Not Null,
    Created_At    Timestamptz      Default Current_Timestamp,
    Updated_At    Timestamptz      Default Current_Timestamp,
    Cancelled_At  Timestamptz      Default Null
);

Create Table Sales_Order_Detail
(
    Id             Uuid Primary Key Default Uuid_Generate_V4(),
    Sales_Order_Id Uuid References Sales_Order (Id) On Delete Cascade,
    Product_Id     Uuid References Product (Id) On Delete Cascade,
    Quantity       INT            Not Null,
    Price          NUMERIC(15, 2) Not Null,
    Created_At     Timestamptz      Default Current_Timestamp,
    Updated_At     Timestamptz      Default Current_Timestamp,
    Deleted_At     Timestamptz      Default Null
);

Create Table Delivery_Note
(
    Id             Uuid Primary Key      Default Uuid_Generate_V4(),
    Sales_Order_Id Uuid References Sales_Order (Id) On Delete Cascade,
    Delivery_Date  Timestamptz           Default Current_Timestamp,
    Recipient_Name VARCHAR(255) Not Null,
    Is_Received    BOOLEAN      Not Null Default False,
    Created_At     Timestamptz           Default Current_Timestamp,
    Updated_At     Timestamptz           Default Current_Timestamp,
    Cancelled_At   Timestamptz           Default Null
);

Create Table Customer
(
    Id         Uuid Primary Key Default Uuid_Generate_V4(),
    Name       VARCHAR(255) Not Null,
    Address    TEXT,
    Telephone  VARCHAR(50),
    Created_At Timestamptz
);

-- Indexes to improve query performance
Create Index Idx_Financial_Transactions_User_Id On Financial_Transaction (User_Id);
Create Index Idx_Users_Username On Appuser (Username);
Create Index Idx_Employees_Name On Employee (Name);
Create Index Idx_Employees_Position On Employee (Position);
Create Index Idx_Wages_Employee_Id On Wage (Employee_Id);
Create Index Idx_Wages_Period On Wage (Month, Year);
Create Index Idx_Wage_Details_Wage_Id On Wage_Detail (Wage_Id);
Create Index Idx_Wage_Details_Component_Name On Wage_Detail (Component_Name);
Create Index Idx_Employees_Nik On Employee (Nik);
Create Index Idx_Sales_Orders_Status On Sales_Order (Status);
Create Index Idx_Products_Name On Product (Name);
Create Index Idx_Storages_Name On Storage (Name);
Create Unique Index Idx_Inventory_Product_Storage On Inventory (Product_Id, Storage_Id);
Create Index Idx_Inventory_Quantity On Inventory (Quantity);
Create Index Idx_Inventory_Logs_Inventory_Id On Inventory_Log (Inventory_Id);
Create Index Idx_Inventory_Logs_Action On Inventory_Log (Action);