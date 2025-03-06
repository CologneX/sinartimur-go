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

-- Table: HR Management
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
    Name        VARCHAR(255) Not Null,
    Description TEXT,
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

Create Table Customer
(
    Id         Uuid Primary Key Default Uuid_Generate_V4(),
    Name       VARCHAR(255) Not Null,
    Address    TEXT,
    Telephone  VARCHAR(50),
    Created_At Timestamptz      Default Current_Timestamp,
    Updated_At Timestamptz      Default Current_Timestamp,
    Deleted_At Timestamptz      Default Null
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
    Id                  Uuid Primary Key Default Uuid_Generate_V4(),
    Supplier_Id         Uuid             Default Null References Supplier (Id) On Delete Set Null,
    Order_Date          Timestamptz      Default Current_Timestamp,
    Status              VARCHAR(50)                   Not Null, -- ordered, received, checked, completed, partially_returned, returned, cancelled
    Total_Amount        NUMERIC(15, 2)                Not Null,
    Payment_Due_Date    Timestamptz      Default Null,
    Created_By          Uuid                          Not Null References Appuser (Id) On Delete Set Null,
    Received_By         Uuid                          References Appuser (Id) On Delete Set Null,
    Checked_By          Uuid                          References Appuser (Id) On Delete Set Null,
    Fully_Returned_By   Uuid                          References Appuser (Id) On Delete Set Null,
    Return_Cancelled_By Uuid                          References Appuser (Id) On Delete Set Null,
    Cancelled_By        Uuid                          References Appuser (Id) On Delete Set Null,
    Created_At          Timestamptz      Default Current_Timestamp,
    Updated_At          Timestamptz      Default Current_Timestamp,
    Cancelled_At        Timestamptz      Default Null,
    Received_At         Timestamptz      Default Null,
    Checked_At          Timestamptz      Default Null,
    Fully_Returned_At   Timestamptz      Default Null,
    Return_Cancelled_At Timestamptz      Default Null
);

Create Table Purchase_Order_Detail
(
    Id                      Uuid Primary Key Default Uuid_Generate_V4(),
    Purchase_Order_Id       Uuid References Purchase_Order (Id) On Delete Cascade,
    Product_Id              Uuid References Product (Id) On Delete Cascade,
    Requested_Quantity      NUMERIC(15, 2) Not Null,
    Total_Returned_Quantity NUMERIC(15, 2)   Default 0,
    Received_Quantity       NUMERIC(15, 2)   Default 0,
    Unit_Price              NUMERIC(15, 2) Not Null,
    Created_At              Timestamptz      Default Current_Timestamp,
    Updated_At              Timestamptz      Default Current_Timestamp
);

-- Product Batch table to track products by SKU
Create Table Product_Batch
(
    Id                Uuid Primary Key Default Uuid_Generate_V4(),
    Sku               VARCHAR(100)   Not Null Unique,
    Product_Id        Uuid           Not Null References Product (Id) On Delete Cascade,
    Purchase_Order_Id Uuid           Not Null References Purchase_Order (Id) On Delete Cascade,
    Initial_Quantity  NUMERIC(15, 2) Not Null,
    Current_Quantity  NUMERIC(15, 2) Not Null,
    Unit_Price        NUMERIC(15, 2) Not Null,
    Created_At        Timestamptz      Default Current_Timestamp,
    Updated_At        Timestamptz      Default Current_Timestamp
);

-- Batch Storage table to track where batches are stored
Create Table Batch_Storage
(
    Id         Uuid Primary Key        Default Uuid_Generate_V4(),
    Batch_Id   Uuid References Product_Batch (Id) On Delete Cascade,
    Storage_Id Uuid References Storage (Id) On Delete Cascade,
    Quantity   NUMERIC(15, 2) Not Null Default 0,
    Created_At Timestamptz             Default Current_Timestamp,
    Updated_At Timestamptz             Default Current_Timestamp,
    Unique (Batch_Id, Storage_Id)
);

Create Table Purchase_Order_Return
(
    Id                 Uuid Primary Key        Default Uuid_Generate_V4(),
    Purchase_Order_Id  Uuid References Purchase_Order (Id) On Delete Cascade,
    Product_Detail_Id  Uuid References Purchase_Order_Detail (Id) On Delete Cascade,
    Return_Quantity    NUMERIC(15, 2) Not Null,
    Remaining_Quantity NUMERIC(15, 2) Not Null,
    Return_Reason      TEXT                    Default Null,
    Return_Status      VARCHAR(50)    Not Null, -- pending, completed, cancelled
    Returned_By        Uuid           References Appuser (Id) On Delete Set Null,
    Cancelled_By       Uuid           References Appuser (Id) On Delete Set Null,
    Returned_At        Timestamptz    Not Null Default Current_Timestamp,
    Cancelled_At       Timestamptz             Default Null
);

Create Table Purchase_Order_Return_Batch
(
    Id                 Uuid Primary Key Default Uuid_Generate_V4(),
    Purchase_Return_Id Uuid References Purchase_Order_Return (Id) On Delete Cascade,
    Batch_Id           Uuid           References Product_Batch (Id) On Delete Set Null,
    Return_Quantity    NUMERIC(15, 2) Not Null,
    Created_At         Timestamptz      Default Current_Timestamp
);

-- Table: Sales
Create Table Sales_Order
(
    Id                    Uuid Primary Key Default Uuid_Generate_V4(),
    Customer_Id           Uuid           References Customer (Id) On Delete Set Null,
    Customer_Name         VARCHAR(255)   Not Null,
    Order_Date            Timestamptz      Default Current_Timestamp,
    Status                VARCHAR(50)    Not Null, -- draft, invoiced, delivered, partially_returned, returned, cancelled
    Payment_Method        VARCHAR(50)    Not Null, -- cash, paylater
    Payment_Due_Date      Timestamptz      Default Null,
    Total_Amount          NUMERIC(15, 2) Not Null,
    Invoice_Created       BOOLEAN          Default False,
    Delivery_Note_Created BOOLEAN          Default False,
    Created_By            Uuid           Not Null References Appuser (Id) On Delete Set Null,
    Invoiced_By           Uuid           References Appuser (Id) On Delete Set Null,
    Delivered_By          Uuid           References Appuser (Id) On Delete Set Null,
    Cancelled_By          Uuid           References Appuser (Id) On Delete Set Null,
    Fully_Returned_By     Uuid           References Appuser (Id) On Delete Set Null,
    Return_Cancelled_By   Uuid           References Appuser (Id) On Delete Set Null,
    Created_At            Timestamptz      Default Current_Timestamp,
    Updated_At            Timestamptz      Default Current_Timestamp,
    Invoiced_At           Timestamptz      Default Null,
    Delivered_At          Timestamptz      Default Null,
    Cancelled_At          Timestamptz      Default Null,
    Fully_Returned_At     Timestamptz      Default Null,
    Return_Cancelled_At   Timestamptz      Default Null
);

Create Table Sales_Order_Detail
(
    Id             Uuid Primary Key Default Uuid_Generate_V4(),
    Sales_Order_Id Uuid References Sales_Order (Id) On Delete Cascade,
    Batch_Id       Uuid References Product_Batch (Id) On Delete Cascade,
    Product_Id     Uuid References Product (Id) On Delete Cascade,
    Quantity       NUMERIC(15, 2) Not Null,
    Unit_Price     NUMERIC(15, 2) Not Null,
    Created_At     Timestamptz      Default Current_Timestamp,
    Updated_At     Timestamptz      Default Current_Timestamp,
    Deleted_At     Timestamptz      Default Null
);

-- Sales_Order_Storage table to track which storage items are taken from
Create Table Sales_Order_Storage
(
    Id                    Uuid Primary Key Default Uuid_Generate_V4(),
    Sales_Order_Detail_Id Uuid References Sales_Order_Detail (Id) On Delete Cascade,
    Storage_Id            Uuid References Storage (Id) On Delete Cascade,
    Batch_Id              Uuid References Product_Batch (Id) On Delete Cascade,
    Quantity              NUMERIC(15, 2) Not Null,
    Created_At            Timestamptz      Default Current_Timestamp,
    Updated_At            Timestamptz      Default Current_Timestamp
);

Create Table Sales_Order_Return
(
    Id                      Uuid Primary Key        Default Uuid_Generate_V4(),
    Sales_Order_Id          Uuid References Sales_Order (Id) On Delete Cascade,
    Sales_Detail_Id         Uuid References Sales_Order_Detail (Id) On Delete Cascade,
    Return_Quantity         NUMERIC(15, 2) Not Null,
    Remaining_Quantity      NUMERIC(15, 2) Not Null,
    Total_Returned_Quantity NUMERIC(15, 2)          Default 0,
    Return_Reason           TEXT                    Default Null,
    Return_Status           VARCHAR(50)    Not Null, -- pending, completed, cancelled
    Returned_By             Uuid           References Appuser (Id) On Delete Set Null,
    Cancelled_By            Uuid           References Appuser (Id) On Delete Set Null,
    Returned_At             Timestamptz    Not Null Default Current_Timestamp,
    Cancelled_At            Timestamptz             Default Null
);

Create Table Sales_Order_Return_Batch
(
    Id              Uuid Primary Key Default Uuid_Generate_V4(),
    Sales_Return_Id Uuid References Sales_Order_Return (Id) On Delete Cascade,
    Batch_Id        Uuid           References Product_Batch (Id) On Delete Set Null,
    Return_Quantity NUMERIC(15, 2) Not Null,
    Created_At      Timestamptz      Default Current_Timestamp
);

Create Table Delivery_Note
(
    Id             Uuid Primary Key Default Uuid_Generate_V4(),
    Sales_Order_Id Uuid References Sales_Order (Id) On Delete Cascade,
    Delivery_Date  Timestamptz      Default Current_Timestamp,
    Driver_Name    VARCHAR(255) Not Null,
    Recipient_Name VARCHAR(255) Not Null,
    Created_By     Uuid         Not Null References Appuser (Id) On Delete Set Null,
    Created_At     Timestamptz      Default Current_Timestamp,
    Updated_At     Timestamptz      Default Current_Timestamp,
    Cancelled_At   Timestamptz      Default Null,
    Cancelled_By   Uuid         References Appuser (Id) On Delete Set Null
);

Create Table Inventory_Log
(
    Id                Uuid Primary Key Default Uuid_Generate_V4(),
    Batch_Id          Uuid           References Product_Batch (Id) On Delete Set Null,
    Storage_Id        Uuid           References Storage (Id) On Delete Set Null,
    User_Id           Uuid References Appuser (Id),
    Purchase_Order_Id Uuid           References Purchase_Order (Id) On Delete Set Null,
    Sales_Order_Id    Uuid           References Sales_Order (Id) On Delete Set Null,
    Action            VARCHAR(50)    Not Null, -- e.g., "add", "remove", "transfer", "return"
    Quantity          NUMERIC(15, 2) Not Null,
    Log_Date          Timestamptz      Default Current_Timestamp,
    Description       TEXT,
    Created_At        Timestamptz      Default Current_Timestamp,
    Updated_At        Timestamptz      Default Current_Timestamp
);

-- Table: Financial Transactions
Create Table Financial_Transaction
(
    Id                Uuid Primary Key Default Uuid_Generate_V4(),
    User_Id           Uuid References Appuser (Id),
    Amount            NUMERIC(15, 2) Not Null,
    Type              VARCHAR(50)    Not Null,
    Purchase_Order_Id Uuid           References Purchase_Order (Id) On Delete Set Null,
    Sales_Order_Id    Uuid           References Sales_Order (Id) On Delete Set Null,
    Description       TEXT,
    Transaction_Date  Timestamptz    Not Null,
    Created_At        Timestamptz      Default Current_Timestamp,
    Edited_At         Timestamptz      Default Null,
    Deleted_At        Timestamptz      Default Null
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
Create Index Idx_Inventory_Logs_Log_Date On Inventory_Log (Log_Date);
Create Index Idx_Inventory_Logs_Action On Inventory_Log (Action);
Create Index Idx_Product_Batch_Sku On Product_Batch (Sku);
Create Index Idx_Product_Batch_Product_Id On Product_Batch (Product_Id);
Create Index Idx_Batch_Storage_Batch_Id On Batch_Storage (Batch_Id);
Create Index Idx_Sales_Order_Payment_Method On Sales_Order (Payment_Method);
Create Index Idx_Sales_Order_Detail_Batch_Id On Sales_Order_Detail (Batch_Id);
Create Index Idx_Sales_Order_Storage_Batch_Id On Sales_Order_Storage (Batch_Id);
Create Index Idx_Inventory_Log_Batch_Id On Inventory_Log (Batch_Id);
Create Index Idx_Purchase_Order_Created_By On Purchase_Order (Created_By);
Create Index Idx_Purchase_Order_Received_By On Purchase_Order (Received_By);
Create Index Idx_Purchase_Order_Checked_By On Purchase_Order (Checked_By);
Create Index Idx_Purchase_Order_Cancelled_By On Purchase_Order (Cancelled_By);
Create Index Idx_Sales_Order_Created_By On Sales_Order (Created_By);
Create Index Idx_Sales_Order_Invoiced_By On Sales_Order (Invoiced_By);
Create Index Idx_Sales_Order_Delivered_By On Sales_Order (Delivered_By);
Create Index Idx_Sales_Order_Cancelled_By On Sales_Order (Cancelled_By);
Create Index Idx_Delivery_Note_Created_By On Delivery_Note (Created_By);
Create Index Idx_Delivery_Note_Cancelled_By On Delivery_Note (Cancelled_By);
Create Index Idx_Purchase_Order_Return_Order_Id On Purchase_Order_Return (Purchase_Order_Id);
Create Index Idx_Purchase_Order_Return_Detail_Id On Purchase_Order_Return (Product_Detail_Id);
Create Index Idx_Purchase_Order_Return_Batch_Return_Id On Purchase_Order_Return_Batch (Purchase_Return_Id);
Create Index Idx_Purchase_Order_Return_Returned_By On Purchase_Order_Return (Returned_By);
Create Index Idx_Purchase_Order_Fully_Returned_By On Purchase_Order (Fully_Returned_By);
Create Index Idx_Sales_Order_Fully_Returned_By On Sales_Order (Fully_Returned_By);
Create Index Idx_Sales_Order_Return_Cancelled_By On Sales_Order (Return_Cancelled_By);
Create Index Idx_Sales_Order_Return_Order_Id On Sales_Order_Return (Sales_Order_Id);
Create Index Idx_Sales_Order_Return_Detail_Id On Sales_Order_Return (Sales_Detail_Id);
Create Index Idx_Sales_Order_Return_Batch_Return_Id On Sales_Order_Return_Batch (Sales_Return_Id);