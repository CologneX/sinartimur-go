-- ================================================
-- Seed Data for HR Management
-- ================================================

-- Insert Employees
INSERT INTO Employee (Id, Name, Position, Phone, Nik, Hired_Date)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'John Doe', 'Developer', '1234567890', 'NIK001', CURRENT_TIMESTAMP),
    ('22222222-2222-2222-2222-222222222222', 'Jane Smith', 'Manager', '0987654321', 'NIK002', CURRENT_TIMESTAMP);

-- Insert Wages for Employees
INSERT INTO Wage (Id, Employee_Id, Total_Amount, Month, Year)
VALUES
    ('33333333-3333-3333-3333-333333333333', '11111111-1111-1111-1111-111111111111', 5000.00, 8, 2025),
    ('44444444-4444-4444-4444-444444444444', '22222222-2222-2222-2222-222222222222', 7000.00, 8, 2025);

-- Insert Wage Details
INSERT INTO Wage_Detail (Id, Wage_Id, Component_Name, Description, Amount)
VALUES
    ('55555555-5555-5555-5555-555555555555', '33333333-3333-3333-3333-333333333333', 'Base Salary', 'Monthly base salary', 4000.00),
    ('66666666-6666-6666-6666-666666666666', '33333333-3333-3333-3333-333333333333', 'Bonus', 'Performance bonus', 1000.00),
    ('77777777-7777-7777-7777-777777777777', '44444444-4444-4444-4444-444444444444', 'Base Salary', 'Monthly base salary', 7000.00);

-- ================================================
-- Seed Data for Inventory
-- ================================================

-- Insert Categories
INSERT INTO Category (Id, Name, Description)
VALUES
    ('88888888-8888-8888-8888-888888888888', 'Electronics', 'Electronic devices'),
    ('99999999-9999-9999-9999-999999999999', 'Furniture', 'Office furniture');

-- Insert Units
INSERT INTO Unit (Id, Name, Description)
VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Piece', 'Single piece'),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Box', 'Packaged box');

-- Insert Products
INSERT INTO Product (Id, Name, Description, Category_Id, Unit_Id)
VALUES
    ('c0c0c0c0-c0c0-c0c0-c0c0-c0c0c0c0c0c0', 'Laptop', 'Gaming laptop', '88888888-8888-8888-8888-888888888888', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    ('d0d0d0d0-d0d0-d0d0-d0d0-d0d0d0d0d0d0', 'Desk', 'Office desk', '99999999-9999-9999-9999-999999999999', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa');

-- Insert Storages
INSERT INTO Storage (Id, Name, Location)
VALUES
    ('e0e0e0e0-e0e0-e0e0-e0e0-e0e0e0e0e0e0', 'Main Warehouse', 'Building A'),
    ('f0f0f0f0-f0f0-f0f0-f0f0-f0f0f0f0f0f0', 'Secondary Warehouse', 'Building B');

-- Insert Customers
INSERT INTO Customer (Id, Name, Address, Telephone)
VALUES
    ('01010101-0101-0101-0101-010101010101', 'Acme Corp', '123 Business Rd', '111-222-3333'),
    ('02020202-0202-0202-0202-020202020202', 'Beta LLC', '456 Commerce St', '444-555-6666');

-- Insert Suppliers
INSERT INTO Supplier (Id, Name, Address, Telephone)
VALUES
    ('03030303-0303-0303-0303-030303030303', 'Tech Supplies Inc.', '789 Industry Ave', '777-888-9999'),
    ('04040404-0404-0404-0404-040404040404', 'Furniture Co.', '321 Home St', '000-111-2222');

-- ================================================
-- Seed Data for Purchase
-- ================================================

-- Insert Purchase Orders (using the provided Appuser for Created_By and other related columns)
INSERT INTO Purchase_Order (Id, Supplier_Id, Order_Date, Status, Total_Amount, Payment_Due_Date, Created_By, Received_By)
VALUES
    ('05050505-0505-0505-0505-050505050505', '03030303-0303-0303-0303-030303030303', CURRENT_TIMESTAMP, 'ordered', 1500.00, CURRENT_TIMESTAMP + INTERVAL '30 days', '22196a3c-a2af-4b32-a89d-0f915b2847f3', NULL),
    ('06060606-0606-0606-0606-060606060606', '04040404-0404-0404-0404-040404040404', CURRENT_TIMESTAMP, 'received', 800.00, CURRENT_TIMESTAMP + INTERVAL '30 days', '22196a3c-a2af-4b32-a89d-0f915b2847f3', '22196a3c-a2af-4b32-a89d-0f915b2847f3');

-- Insert Purchase Order Details
INSERT INTO Purchase_Order_Detail (Id, Purchase_Order_Id, Product_Id, Requested_Quantity, Unit_Price)
VALUES
    ('07070707-0707-0707-0707-070707070707', '05050505-0505-0505-0505-050505050505', 'c0c0c0c0-c0c0-c0c0-c0c0-c0c0c0c0c0c0', 10, 150.00),
    ('08080808-0808-0808-0808-080808080808', '06060606-0606-0606-0606-060606060606', 'd0d0d0d0-d0d0-d0d0-d0d0-d0d0d0d0d0d0', 5, 160.00);

-- Insert Product Batches (track SKU and stock)
INSERT INTO Product_Batch (Id, Sku, Product_Id, Purchase_Order_Id, Initial_Quantity, Current_Quantity, Unit_Price)
VALUES
    ('09090909-0909-0909-0909-090909090909', 'SKU-LAP-001', 'c0c0c0c0-c0c0-c0c0-c0c0-c0c0c0c0c0c0', '05050505-0505-0505-0505-050505050505', 10, 10, 150.00),
    ('0a0a0a0a-0a0a-0a0a-0a0a-0a0a0a0a0a0a', 'SKU-DESK-001', 'd0d0d0d0-d0d0-d0d0-d0d0-d0d0d0d0d0d0', '06060606-0606-0606-0606-060606060606', 5, 5, 160.00);

-- Insert Batch Storage mappings
INSERT INTO Batch_Storage (Id, Batch_Id, Storage_Id, Quantity)
VALUES
    ('0b0b0b0b-0b0b-0b0b-0b0b-0b0b0b0b0b0b', '09090909-0909-0909-0909-090909090909', 'e0e0e0e0-e0e0-e0e0-e0e0-e0e0e0e0e0e0', 10),
    ('0c0c0c0c-0c0c-0c0c-0c0c-0c0c0c0c0c0c', '0a0a0a0a-0a0a-0a0a-0a0a-0a0a0a0a0a0a', 'f0f0f0f0-f0f0-f0f0-f0f0-f0f0f0f0f0f0', 5);

-- Insert a Purchase Order Return for the first purchase order detail
INSERT INTO Purchase_Order_Return (Id, Purchase_Order_Id, Product_Detail_Id, Return_Quantity, Remaining_Quantity, Return_Reason, Return_Status, Returned_By)
VALUES
    ('0d0d0d0d-0d0d-0d0d-0d0d-0d0d0d0d0d0d', '05050505-0505-0505-0505-050505050505', '07070707-0707-0707-0707-070707070707', 2, 8, 'Damaged items', 'pending', '22196a3c-a2af-4b32-a89d-0f915b2847f3');

-- Link the Purchase Return to a specific Batch
INSERT INTO Purchase_Order_Return_Batch (Id, Purchase_Return_Id, Batch_Id, Return_Quantity)
VALUES
    ('0e0e0e0e-0e0e-0e0e-0e0e-0e0e0e0e0e0e', '0d0d0d0d-0d0d-0d0d-0d0d-0d0d0d0d0d0d', '09090909-0909-0909-0909-090909090909', 2);

-- ================================================
-- Seed Data for Sales
-- ================================================

-- Insert Sales Orders
INSERT INTO Sales_Order (Id, Customer_Id, Customer_Name, Order_Date, Status, Payment_Method, Total_Amount, Invoice_Created, Created_By)
VALUES
    ('0f0f0f0f-0f0f-0f0f-0f0f-0f0f0f0f0f0f', '01010101-0101-0101-0101-010101010101', 'Acme Corp', CURRENT_TIMESTAMP, 'invoiced', 'cash', 1500.00, true, '22196a3c-a2af-4b32-a89d-0f915b2847f3'),
    ('10101010-1010-1010-1010-101010101010', '02020202-0202-0202-0202-020202020202', 'Beta LLC', CURRENT_TIMESTAMP, 'draft', 'paylater', 800.00, false, '22196a3c-a2af-4b32-a89d-0f915b2847f3');

-- Insert Sales Order Details
INSERT INTO Sales_Order_Detail (Id, Sales_Order_Id, Batch_Id, Product_Id, Quantity, Unit_Price)
VALUES
    ('11111111-2222-3333-4444-555555555555', '0f0f0f0f-0f0f-0f0f-0f0f-0f0f0f0f0f0f', '09090909-0909-0909-0909-090909090909', 'c0c0c0c0-c0c0-c0c0-c0c0-c0c0c0c0c0c0', 1, 1500.00),
    ('66666666-7777-8888-9999-aaaaaaaaaaaa', '10101010-1010-1010-1010-101010101010', '0a0a0a0a-0a0a-0a0a-0a0a-0a0a0a0a0a0a', 'd0d0d0d0-d0d0-d0d0-d0d0-d0d0d0d0d0d0', 2, 800.00);

-- Insert Sales Order Storage records
INSERT INTO Sales_Order_Storage (Id, Sales_Order_Detail_Id, Storage_Id, Batch_Id, Quantity)
VALUES
    ('12121212-1212-1212-1212-121212121212', '11111111-2222-3333-4444-555555555555', 'e0e0e0e0-e0e0-e0e0-e0e0-e0e0e0e0e0e0', '09090909-0909-0909-0909-090909090909', 1),
    ('13131313-1313-1313-1313-131313131313', '66666666-7777-8888-9999-aaaaaaaaaaaa', 'f0f0f0f0-f0f0-f0f0-f0f0-f0f0f0f0f0f0', '0a0a0a0a-0a0a-0a0a-0a0a-0a0a0a0a0a0a', 2);

-- Insert a Sales Order Return (for the first sales order detail)
INSERT INTO Sales_Order_Return (Id, Sales_Order_Id, Sales_Detail_Id, Return_Quantity, Remaining_Quantity, Total_Returned_Quantity, Return_Reason, Return_Status, Returned_By)
VALUES
    ('14141414-1414-1414-1414-141414141414', '0f0f0f0f-0f0f-0f0f-0f0f-0f0f0f0f0f0f', '11111111-2222-3333-4444-555555555555', 1, 0, 1, 'Customer returned item', 'completed', '22196a3c-a2af-4b32-a89d-0f915b2847f3');

-- Link the Sales Return to a Batch
INSERT INTO Sales_Order_Return_Batch (Id, Sales_Return_Id, Batch_Id, Return_Quantity)
VALUES
    ('15151515-1515-1515-1515-151515151515', '14141414-1414-1414-1414-141414141414', '09090909-0909-0909-0909-090909090909', 1);

-- Insert a Delivery Note for a Sales Order
INSERT INTO Delivery_Note (Id, Sales_Order_Id, Delivery_Date, Driver_Name, Recipient_Name, Created_By)
VALUES
    ('16161616-1616-1616-1616-161616161616', '0f0f0f0f-0f0f-0f0f-0f0f-0f0f0f0f0f0f', CURRENT_TIMESTAMP, 'Alex Driver', 'John Recipient', '22196a3c-a2af-4b32-a89d-0f915b2847f3');

-- ================================================
-- Seed Data for Inventory & Financial Transactions
-- ================================================

-- Insert an Inventory Log entry
INSERT INTO Inventory_Log (Id, Batch_Id, Storage_Id, User_Id, Purchase_Order_Id, Sales_Order_Id, Action, Quantity, Description)
VALUES
    ('17171717-1717-1717-1717-171717171717', '09090909-0909-0909-0909-090909090909', 'e0e0e0e0-e0e0-e0e0-e0e0-e0e0e0e0e0e0', '22196a3c-a2af-4b32-a89d-0f915b2847f3', '05050505-0505-0505-0505-050505050505', '0f0f0f0f-0f0f-0f0f-0f0f-0f0f0f0f0f0f', 'add', 10, 'Initial stock added');

-- Insert a Financial Transaction entry
INSERT INTO Financial_Transaction (Id, User_Id, Amount, Type, Sales_Order_Id, Description, Transaction_Date)
VALUES
    ('18181818-1818-1818-1818-181818181818', '22196a3c-a2af-4b32-a89d-0f915b2847f3', 1500.00, 'sale', '0f0f0f0f-0f0f-0f0f-0f0f-0f0f0f0f0f0f', 'Sale transaction', CURRENT_TIMESTAMP);
