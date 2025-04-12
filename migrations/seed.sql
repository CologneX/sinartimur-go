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
    ('d0d0d0d0-d0d0-d0d0-d0d0-d0d0d0d0d0d0', 'Desk', 'Office desk', '99999999-9999-9999-9999-999999999999', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    ('d1d1d1d1-d1d1-d1d1-d1d1-d1d1d1d1d1d1', 'Monitor', 'LCD Monitor', '88888888-8888-8888-8888-888888888888', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa');

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

-- ================================================
-- Seed Data for Purchase
-- ================================================

-- Initialize PO document counter
INSERT INTO Document_Counter (Document_Type, Year, Month, Day, Counter)
VALUES
    (
        'PO',
        EXTRACT(YEAR FROM CURRENT_DATE),
        EXTRACT(MONTH FROM CURRENT_DATE),
        EXTRACT(DAY FROM CURRENT_DATE),
        1
    )
ON CONFLICT (Document_Type, Year, Month, Day) 
DO UPDATE SET Counter = Document_Counter.Counter;

-- Insert Suppliers
INSERT INTO Supplier (Id, Name, Address, Telephone)
VALUES
    ('03030303-0303-0303-0303-030303030303', 'Alpha Omega Sigma', '789 Industry Ave', '777-888-9999'),
    ('04040404-0404-0404-0404-040404040404', 'Furniture Masters', '321 Home St', '000-111-2222');

-- Insert Purchase Orders with different statuses according to business logic
INSERT INTO Purchase_Order (
    Id, 
    Serial_Id, 
    Supplier_Id, 
    Order_Date, 
    Status, 
    Total_Amount, 
    Payment_Method,
    Payment_Due_Date, 
    Created_By, 
    Checked_By
)
VALUES
    -- Ordered status
    ('05050505-0505-0505-0505-050505050505', 
     'PO-20250413-0001', 
     '03030303-0303-0303-0303-030303030303', 
     CURRENT_TIMESTAMP, 
     'ordered', 
     1500.00, 
     'cash',
     NULL, 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     NULL),
     
    -- Completed status
    ('06060606-0606-0606-0606-060606060606', 
     'PO-20250413-0002', 
     '04040404-0404-0404-0404-040404040404', 
     CURRENT_TIMESTAMP - INTERVAL '1 day', 
     'completed', 
     800.00, 
     'credit',
     CURRENT_TIMESTAMP + INTERVAL '30 days', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9'),
     
    -- Partially_returned status
    ('07070707-0707-0707-0707-070707070707', 
     'PO-20250413-0003', 
     '03030303-0303-0303-0303-030303030303', 
     CURRENT_TIMESTAMP - INTERVAL '2 days', 
     'partially_returned', 
     2500.00, 
     'cash',
     NULL, 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9'),
     
    -- Returned status
    ('08080808-0808-0808-0808-080808080808', 
     'PO-20250413-0004', 
     '04040404-0404-0404-0404-040404040404', 
     CURRENT_TIMESTAMP - INTERVAL '3 days', 
     'returned', 
     750.00, 
     'cash',
     NULL, 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9'),
     
    -- Cancelled status
    ('09090909-0909-0909-0909-090909090909', 
     'PO-20250413-0005', 
     '03030303-0303-0303-0303-030303030303', 
     CURRENT_TIMESTAMP - INTERVAL '4 days', 
     'cancelled', 
     1800.00, 
     'credit',
     CURRENT_TIMESTAMP + INTERVAL '45 days', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     NULL);

-- Update cancelled PO with cancelled info
UPDATE Purchase_Order
SET 
    Cancelled_By = '20554914-c187-4d8d-a2d1-e5fe4db272e9',
    Cancelled_At = CURRENT_TIMESTAMP - INTERVAL '4 days'
WHERE Id = '09090909-0909-0909-0909-090909090909';

-- Insert Purchase Order Details
INSERT INTO Purchase_Order_Detail (
    Id, 
    Purchase_Order_Id, 
    Product_Id, 
    Requested_Quantity, 
    Unit_Price
)
VALUES
    -- For Ordered PO (05050505-0505-0505-0505-050505050505)
    ('a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 
     '05050505-0505-0505-0505-050505050505', 
     'c0c0c0c0-c0c0-c0c0-c0c0-c0c0c0c0c0c0', 
     5, 
     300.00),
     
    -- For Completed PO (06060606-0606-0606-0606-060606060606)
    ('a2a2a2a2-a2a2-a2a2-a2a2-a2a2a2a2a2a2', 
     '06060606-0606-0606-0606-060606060606', 
     'd0d0d0d0-d0d0-d0d0-d0d0-d0d0d0d0d0d0', 
     4, 
     200.00),
     
    -- For Partially_returned PO (07070707-0707-0707-0707-070707070707)
    ('a3a3a3a3-a3a3-a3a3-a3a3-a3a3a3a3a3a3', 
     '07070707-0707-0707-0707-070707070707', 
     'c0c0c0c0-c0c0-c0c0-c0c0-c0c0c0c0c0c0', 
     5, 
     300.00),
    ('a4a4a4a4-a4a4-a4a4-a4a4-a4a4a4a4a4a4', 
     '07070707-0707-0707-0707-070707070707', 
     'd1d1d1d1-d1d1-d1d1-d1d1-d1d1d1d1d1d1', 
     10, 
     100.00),
     
    -- For Returned PO (08080808-0808-0808-0808-080808080808)
    ('a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', 
     '08080808-0808-0808-0808-080808080808', 
     'd0d0d0d0-d0d0-d0d0-d0d0-d0d0d0d0d0d0', 
     3, 
     250.00),
     
    -- For Cancelled PO (09090909-0909-0909-0909-090909090909)
    ('a6a6a6a6-a6a6-a6a6-a6a6-a6a6a6a6a6a6', 
     '09090909-0909-0909-0909-090909090909', 
     'c0c0c0c0-c0c0-c0c0-c0c0-c0c0c0c0c0c0', 
     6, 
     300.00);

-- Insert Product Batches (only for completed or partially_returned POs)
INSERT INTO Product_Batch (
    Id, 
    Sku, 
    Product_Id, 
    Purchase_Order_Id, 
    Initial_Quantity, 
    Current_Quantity, 
    Unit_Price
)
VALUES
    -- For Completed PO (06060606-0606-0606-0606-060606060606)
    ('b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1b1', 
     'DES-FM130425', 
     'd0d0d0d0-d0d0-d0d0-d0d0-d0d0d0d0d0d0', 
     '06060606-0606-0606-0606-060606060606', 
     4, 
     4, 
     200.00),
     
    -- For Partially_returned PO (07070707-0707-0707-0707-070707070707)
    ('b2b2b2b2-b2b2-b2b2-b2b2-b2b2b2b2b2b2', 
     'LAP-AOS130425', 
     'c0c0c0c0-c0c0-c0c0-c0c0-c0c0c0c0c0c0', 
     '07070707-0707-0707-0707-070707070707', 
     5, 
     3, 
     300.00),
    ('b3b3b3b3-b3b3-b3b3-b3b3-b3b3b3b3b3b3', 
     'MON-AOS130425', 
     'd1d1d1d1-d1d1-d1d1-d1d1-d1d1d1d1d1d1', 
     '07070707-0707-0707-0707-070707070707', 
     10, 
     10, 
     100.00),
     
    -- For Returned PO (08080808-0808-0808-0808-080808080808) - keep the batches but zero quantity
    ('b4b4b4b4-b4b4-b4b4-b4b4-b4b4b4b4b4b4', 
     'DES-FM130425-2', 
     'd0d0d0d0-d0d0-d0d0-d0d0-d0d0d0d0d0d0', 
     '08080808-0808-0808-0808-080808080808', 
     3, 
     0, 
     250.00);

-- Insert Batch Storage mappings
INSERT INTO Batch_Storage (
    Id, 
    Batch_Id, 
    Storage_Id, 
    Quantity, 
    Created_At, 
    Updated_At
)
VALUES
    -- For Completed PO batches
    ('c1c1c1c1-c1c1-c1c1-c1c1-c1c1c1c1c1c1', 
     'b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1b1', 
     'e0e0e0e0-e0e0-e0e0-e0e0-e0e0e0e0e0e0', 
     4, 
     CURRENT_TIMESTAMP - INTERVAL '1 day', 
     CURRENT_TIMESTAMP - INTERVAL '1 day'),
     
    -- For Partially_returned PO batches (some in main, some in secondary)
    ('c2c2c2c2-c2c2-c2c2-c2c2-c2c2c2c2c2c2', 
     'b2b2b2b2-b2b2-b2b2-b2b2-b2b2b2b2b2b2', 
     'e0e0e0e0-e0e0-e0e0-e0e0-e0e0e0e0e0e0', 
     3, 
     CURRENT_TIMESTAMP - INTERVAL '2 days', 
     CURRENT_TIMESTAMP - INTERVAL '2 days'),
    ('c3c3c3c3-c3c3-c3c3-c3c3-c3c3c3c3c3c3', 
     'b3b3b3b3-b3b3-b3b3-b3b3-b3b3b3b3b3b3', 
     'f0f0f0f0-f0f0-f0f0-f0f0-f0f0f0f0f0f0', 
     10, 
     CURRENT_TIMESTAMP - INTERVAL '2 days', 
     CURRENT_TIMESTAMP - INTERVAL '2 days'),
     
    -- For Returned PO batches (zero quantity)
    ('c4c4c4c4-c4c4-c4c4-c4c4-c4c4c4c4c4c4', 
     'b4b4b4b4-b4b4-b4b4-b4b4-b4b4b4b4b4b4', 
     'e0e0e0e0-e0e0-e0e0-e0e0-e0e0e0e0e0e0', 
     0, 
     CURRENT_TIMESTAMP - INTERVAL '3 days', 
     CURRENT_TIMESTAMP - INTERVAL '3 days');

-- Insert Purchase Order Returns
INSERT INTO Purchase_Order_Return (
    Id, 
    Purchase_Order_Id, 
    Product_Detail_Id, 
    Return_Quantity, 
    Reason, 
    Status, 
    Returned_By, 
    Returned_At
)
VALUES
    -- Partial return for the partially_returned PO
    ('d1d1d1d1-d1d1-d1d1-d1d1-d1d1d1d1d1d1', 
     '07070707-0707-0707-0707-070707070707', 
     'a3a3a3a3-a3a3-a3a3-a3a3-a3a3a3a3a3a3', 
     2, 
     'Damaged on arrival', 
     'returned', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     CURRENT_TIMESTAMP - INTERVAL '2 days'),
     
    -- Complete return for the returned PO
    ('d2d2d2d2-d2d2-d2d2-d2d2-d2d2d2d2d2d2', 
     '08080808-0808-0808-0808-080808080808', 
     'a5a5a5a5-a5a5-a5a5-a5a5-a5a5a5a5a5a5', 
     3, 
     'Wrong item delivered', 
     'returned', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     CURRENT_TIMESTAMP - INTERVAL '3 days');

-- Insert Purchase Order Return Batches
INSERT INTO Purchase_Order_Return_Batch (
    Id, 
    Purchase_Return_Id, 
    Batch_Id, 
    Quantity, 
    Created_At
)
VALUES
    -- Link return to specific batches
    ('e1e1e1e1-e1e1-e1e1-e1e1-e1e1e1e1e1e1', 
     'd1d1d1d1-d1d1-d1d1-d1d1-d1d1d1d1d1d1', 
     'b2b2b2b2-b2b2-b2b2-b2b2-b2b2b2b2b2b2', 
     2, 
     CURRENT_TIMESTAMP - INTERVAL '2 days'),
     
    ('e2e2e2e2-e2e2-e2e2-e2e2-e2e2e2e2e2e2', 
     'd2d2d2d2-d2d2-d2d2-d2d2-d2d2d2d2d2d2', 
     'b4b4b4b4-b4b4-b4b4-b4b4-b4b4b4b4b4b4', 
     3, 
     CURRENT_TIMESTAMP - INTERVAL '3 days');

-- ================================================
-- Inventory and Financial Logs
-- ================================================

-- Insert Inventory Logs for received items
INSERT INTO Inventory_Log (
    Id, 
    Batch_Id, 
    Storage_Id, 
    User_Id, 
    Purchase_Order_Id, 
    Action, 
    Quantity, 
    Description, 
    Log_Date
)
VALUES
    -- Receipt of completed PO
    ('f1f1f1f1-f1f1-f1f1-f1f1-f1f1f1f1f1f1', 
     'b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1b1', 
     'e0e0e0e0-e0e0-e0e0-e0e0-e0e0e0e0e0e0', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     '06060606-0606-0606-0606-060606060606', 
     'add', 
     4, 
     'Pembelian Barang PO-20250413-0002', 
     CURRENT_TIMESTAMP - INTERVAL '1 day'),
     
    -- Receipt of partially_returned PO
    ('f2f2f2f2-f2f2-f2f2-f2f2-f2f2f2f2f2f2', 
     'b2b2b2b2-b2b2-b2b2-b2b2-b2b2b2b2b2b2', 
     'e0e0e0e0-e0e0-e0e0-e0e0-e0e0e0e0e0e0', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     '07070707-0707-0707-0707-070707070707', 
     'add', 
     5, 
     'Pembelian Barang PO-20250413-0003', 
     CURRENT_TIMESTAMP - INTERVAL '2 days'),
    ('f3f3f3f3-f3f3-f3f3-f3f3-f3f3f3f3f3f3', 
     'b3b3b3b3-b3b3-b3b3-b3b3-b3b3b3b3b3b3', 
     'f0f0f0f0-f0f0-f0f0-f0f0-f0f0f0f0f0f0', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     '07070707-0707-0707-0707-070707070707', 
     'add', 
     10, 
     'Pembelian Barang PO-20250413-0003', 
     CURRENT_TIMESTAMP - INTERVAL '2 days'),
     
    -- Receipt of fully returned PO
    ('f4f4f4f4-f4f4-f4f4-f4f4-f4f4f4f4f4f4', 
     'b4b4b4b4-b4b4-b4b4-b4b4-b4b4b4b4b4b4', 
     'e0e0e0e0-e0e0-e0e0-e0e0-e0e0e0e0e0e0', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     '08080808-0808-0808-0808-080808080808', 
     'add', 
     3, 
     'Pembelian Barang PO-20250413-0004', 
     CURRENT_TIMESTAMP - INTERVAL '3 days'),
     
    -- Return for the partially_returned PO
    ('f5f5f5f5-f5f5-f5f5-f5f5-f5f5f5f5f5f5', 
     'b2b2b2b2-b2b2-b2b2-b2b2-b2b2b2b2b2b2', 
     'e0e0e0e0-e0e0-e0e0-e0e0-e0e0e0e0e0e0', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     '07070707-0707-0707-0707-070707070707', 
     'return', 
     2, 
     'Return Pembelian PO-20250413-0003', 
     CURRENT_TIMESTAMP - INTERVAL '2 days' + INTERVAL '2 hours'),
     
    -- Return for the fully returned PO
    ('f6f6f6f6-f6f6-f6f6-f6f6-f6f6f6f6f6f6', 
     'b4b4b4b4-b4b4-b4b4-b4b4-b4b4b4b4b4b4', 
     'e0e0e0e0-e0e0-e0e0-e0e0-e0e0e0e0e0e0', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     '08080808-0808-0808-0808-080808080808', 
     'return', 
     3, 
     'Return Pembelian PO-20250413-0004', 
     CURRENT_TIMESTAMP - INTERVAL '3 days' + INTERVAL '3 hours');

-- Insert Financial Transaction Logs
INSERT INTO Financial_Transaction_Log (
    Id, 
    User_Id, 
    Amount, 
    Type, 
    Purchase_Order_Id, 
    Description, 
    Transaction_Date
)
VALUES
    -- Financial transaction for completed PO
    ('c1a1b2c3-d4e5-f6a7-b8c9-d0e1f2a3b4c5', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     800.00, 
     'purchase', 
     '06060606-0606-0606-0606-060606060606', 
     'Pembelian Barang PO-20250413-0002', 
     CURRENT_TIMESTAMP - INTERVAL '1 day'),
     
    -- Financial transaction for partially_returned PO
    ('d2b3c4d5-e6f7-a8b9-c0d1-e2f3a4b5c6d7', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     2500.00, 
     'purchase', 
     '07070707-0707-0707-0707-070707070707', 
     'Pembelian Barang PO-20250413-0003', 
     CURRENT_TIMESTAMP - INTERVAL '2 days'),
     
    -- Financial transaction for fully returned PO
    ('e3c4d5e6-f7a8-b9c0-d1e2-f3a4b5c6d7e8', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     750.00, 
     'purchase', 
     '08080808-0808-0808-0808-080808080808', 
     'Pembelian Barang PO-20250413-0004', 
     CURRENT_TIMESTAMP - INTERVAL '3 days'),
     
    -- Financial transaction for partial return
    ('f4d5e6f7-a8b9-c0d1-e2f3-a4b5c6d7e8f9', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     -600.00, 
     'purchase_return', 
     '07070707-0707-0707-0707-070707070707', 
     'Return Pembelian PO-20250413-0003', 
     CURRENT_TIMESTAMP - INTERVAL '2 days' + INTERVAL '2 hours'),
     
    -- Financial transaction for full return
    ('a5e6f7a8-b9c0-d1e2-f3a4-b5c6d7e8f9a0', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     -750.00, 
     'purchase_return', 
     '08080808-0808-0808-0808-080808080808', 
     'Return Pembelian PO-20250413-0004', 
     CURRENT_TIMESTAMP - INTERVAL '3 days' + INTERVAL '3 hours'),
     
    -- Financial transaction for cancelled PO
    ('b6f7a8b9-c0d1-e2f3-a4b5-c6d7e8f9a0b1', 
     '20554914-c187-4d8d-a2d1-e5fe4db272e9', 
     0.00, 
     'purchase_cancel', 
     '09090909-0909-0909-0909-090909090909', 
     'Pembatalan Pesanan Pembelian PO-20250413-0005', 
     CURRENT_TIMESTAMP - INTERVAL '4 days');

-- ================================================
-- Seed Data for Sales
-- ================================================

-- Insert Sales Orders
INSERT INTO Sales_Order (Id, serial_id, Customer_Id, Order_Date, Status, Payment_Method, Payment_Due_Date, Total_Amount, Created_By)
VALUES
    ('0f0f0f0f-0f0f-0f0f-0f0f-0f0f0f0f0f0f', 'SO-00001','01010101-0101-0101-0101-010101010101', CURRENT_TIMESTAMP, 'purchase-order', 'cash', NULL, 1500.00, '20554914-c187-4d8d-a2d1-e5fe4db272e9'),
    ('10101010-1010-1010-1010-101010101010', 'SO-00002','02020202-0202-0202-0202-020202020202', CURRENT_TIMESTAMP, 'purchase-order', 'paylater', CURRENT_TIMESTAMP + INTERVAL '30 days', 800.00, '20554914-c187-4d8d-a2d1-e5fe4db272e9');

-- Insert Sales Order Details
INSERT INTO Sales_Order_Detail (Id, Sales_Order_Id, Batch_Storage_Id, Quantity, Unit_Price)
VALUES
    ('11111111-2222-3333-4444-555555555555', '0f0f0f0f-0f0f-0f0f-0f0f-0f0f0f0f0f0f', 'c1c1c1c1-c1c1-c1c1-c1c1-c1c1c1c1c1c1', 1, 300.00),
    ('66666666-7777-8888-9999-aaaaaaaaaaaa', '10101010-1010-1010-1010-101010101010', 'c2c2c2c2-c2c2-c2c2-c2c2-c2c2c2c2c2c2', 2, 400.00);

-- Insert a Sales Invoice
INSERT INTO Sales_Invoice (Id, Sales_Order_Id, Serial_Id, Invoice_Date, Total_Amount, Created_By)
VALUES
    ('19191919-1919-1919-1919-191919191919', '0f0f0f0f-0f0f-0f0f-0f0f-0f0f0f0f0f0f', 'SI-00001', CURRENT_TIMESTAMP, 1500.00, '20554914-c187-4d8d-a2d1-e5fe4db272e9');
