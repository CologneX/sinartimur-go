-- Enable UUID Extension
Create Extension If Not Exists "uuid-ossp";

-- Table: Admin
Create Table Users
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
    IsPurchase   BOOLEAN          Default False,
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
Create Table Employees
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


Create Table Wages
(
    Id           Uuid Primary Key Default Uuid_Generate_V4(),
    Employee_Id  Uuid References Employees (Id) On Delete Cascade,
    Total_Amount NUMERIC(12, 2) Not Null,
    Month        INT            Not Null,
    Year         INT            Not Null,
    Created_At   Timestamptz      Default Current_Timestamp,
    Updated_At   Timestamptz      Default Current_Timestamp,
    Deleted_At   Timestamptz      Default Null
);

Create Table Wage_Details
(
    Id             Uuid Primary Key Default Uuid_Generate_V4(),
    Wage_Id        Uuid References Wages (Id) On Delete Cascade,
    Component_Name VARCHAR(100)   Not Null,
    Description    TEXT,
    Amount         NUMERIC(12, 2) Not Null,
    Created_At     Timestamptz      Default Current_Timestamp,
    Updated_At     Timestamptz      Default Current_Timestamp,
    Deleted_At     Timestamptz      Default Null
);

-- Table: Financial Transactions
Create Table Financial_Transactions
(
    Id               Uuid Primary Key Default Uuid_Generate_V4(),
    User_Id          Uuid References Users (Id),
    Amount           NUMERIC(15, 2) Not Null,
    Type             VARCHAR(50)    Not Null,
    Description      TEXT,
    Transaction_Date Timestamptz    Not Null,
    Created_At       Timestamptz      Default Current_Timestamp,
    Edited_At        Timestamptz      Default Null,
    Deleted_At       Timestamptz      Default Null
);

-- Table: Inventory
Create Table Products
(
    Id          Uuid Primary Key Default Uuid_Generate_V4(),
    Name        VARCHAR(255)       Not Null,
    Sku         VARCHAR(50) Unique Not Null, -- Kode unik produk
    Description TEXT,
    Price       NUMERIC(15, 2)     Not Null,
    Created_At  Timestamptz      Default Current_Timestamp,
    Updated_At  Timestamptz      Default Current_Timestamp,
    Deleted_At  Timestamptz      Default Null
);

Create Table Storages
(
    Id         Uuid Primary Key Default Uuid_Generate_V4(),
    Name       VARCHAR(255) Not Null, -- Nama gudang, e.g., "Gudang Utama", "Gudang Cabang A"
    Location   TEXT         Not Null, -- Lokasi fisik gudang
    Created_At Timestamptz      Default Current_Timestamp,
    Updated_At Timestamptz      Default Current_Timestamp
);

Create Table Inventory
(
    Id               Uuid Primary Key Default Uuid_Generate_V4(),
    Product_Id       Uuid References Products (Id) On Delete Cascade,
    Storage_Id       Uuid References Storages (Id) On Delete Cascade,
    Quantity         INT Not Null     Default 0,
    Minimum_Quantity INT Not Null     Default 0,
    Created_At       Timestamptz      Default Current_Timestamp,
    Updated_At       Timestamptz      Default Current_Timestamp
);


Create Table Inventory_Logs
(
    Id           Uuid Primary Key Default Uuid_Generate_V4(),
    Inventory_Id Uuid References Inventory (Id) On Delete Cascade,
    User_Id      Uuid References Users (Id),
    Action       VARCHAR(50) Not Null, -- e.g., "add", "remove", "transfer"
    Quantity     INT         Not Null,
    Log_Date     Timestamptz      Default Current_Timestamp,
    Description  TEXT,
    Created_At   Timestamptz      Default Current_Timestamp,
    Updated_At   Timestamptz      Default Current_Timestamp
);

-- Table: Sales
Create Table Sales_Orders
(
    Id            Uuid Primary Key Default Uuid_Generate_V4(),
    Customer_Name VARCHAR(255)   Not Null,
    Order_Date    Timestamptz      Default Current_Timestamp,
    Status        VARCHAR(50)    Not Null, -- e.g., pending, confirmed, shipped
    Total_Amount  NUMERIC(15, 2) Not Null,
    Created_At    Timestamptz      Default Current_Timestamp,
    Updated_At    Timestamptz      Default Current_Timestamp,
    Cancelled_At  Timestamptz      Default Null
);


Create Table Sales_Order_Items
(
    Id             Uuid Primary Key Default Uuid_Generate_V4(),
    Sales_Order_Id Uuid References Sales_Orders (Id) On Delete Cascade,
    Inventory_Id   Uuid References Inventory (Id),
    Quantity       INT            Not Null,
    Price          NUMERIC(15, 2) Not Null,
    Created_At     Timestamptz      Default Current_Timestamp,
    Updated_At     Timestamptz      Default Current_Timestamp,
    Cancelled_At   Timestamptz      Default Null
);


Create Table Sales_Invoices
(
    Id             Uuid Primary Key Default Uuid_Generate_V4(),
    Sales_Order_Id Uuid References Sales_Orders (Id) On Delete Cascade,
    Invoice_Date   Timestamptz      Default Current_Timestamp,
    Due_Date       Timestamptz    Not Null,
    Total_Amount   NUMERIC(15, 2) Not Null,
    Payment_Status VARCHAR(50)      Default 'unpaid', -- e.g., unpaid, paid, overdue
    Created_At     Timestamptz      Default Current_Timestamp,
    Updated_At     Timestamptz      Default Current_Timestamp,
    Cancelled_At   Timestamptz      Default Null
);


Create Table Delivery_Notes
(
    Id             Uuid Primary Key Default Uuid_Generate_V4(),
    Sales_Order_Id Uuid References Sales_Orders (Id) On Delete Cascade,
    Delivery_Date  Timestamptz      Default Current_Timestamp,
    Recipient_Name VARCHAR(255) Not Null,
    Status         VARCHAR(50)      Default 'in_transit', -- e.g., in_transit, delivered, returned
    Created_At     Timestamptz      Default Current_Timestamp,
    Updated_At     Timestamptz      Default Current_Timestamp,
    Cancelled_At   Timestamptz      Default Null
);

-- Indexes to improve query performance
Create Index Idx_Financial_Transactions_User_Id On Financial_Transactions (User_Id);
Create Index Idx_Orders_Status On Orders (Status);
Create Index Idx_Order_Items_Order_Id On Order_Items (Order_Id);
Create Index Idx_Users_Username On Users (Username);
Create Index Idx_User_Roles_User_Id On User_Roles (User_Id);
Create Index Idx_Employees_Name On Employees (Name);
Create Index Idx_Employees_Position On Employees (Position);
Create Index Idx_User_Roles_Role_Id On User_Roles (Role_Id);
Create Index Idx_Wages_Employee_Id On Wages (Employee_Id);
Create Index Idx_Wages_Period On Wages (Month, Year);
Create Index Idx_Wage_Details_Wage_Id On Wage_Details (Wage_Id);
Create Index Idx_Wage_Details_Component_Name On Wage_Details (Component_Name);
Create Index Idx_Employees_Nik On Employees (Nik);
Create Index Idx_Delivery_Notes_Status On Delivery_Notes (Status);
Create Index Idx_Sales_Invoices_Payment_Status On Sales_Invoices (Payment_Status);
Create Index Idx_Sales_Order_Items_Sales_Order_Id On Sales_Order_Items (Sales_Order_Id);
Create Index Idx_Sales_Orders_Status On Sales_Orders (Status);
Create Index Idx_Products_Sku On Products (Sku);
Create Index Idx_Products_Name On Products (Name);
Create Index Idx_Storages_Name On Storages (Name);
Create Unique Index Idx_Inventory_Product_Storage On Inventory (Product_Id, Storage_Id);
Create Index Idx_Inventory_Quantity On Inventory (Quantity);
Create Index Idx_Inventory_Logs_Inventory_Id On Inventory_Logs (Inventory_Id);
Create Index Idx_Inventory_Logs_Action On Inventory_Logs (Action);