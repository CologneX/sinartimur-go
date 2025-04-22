-- ================================================
-- Seed Data for HR Management
-- ================================================
-- Insert Employees
INSERT INTO
    Employee (Id, Name, Position, Phone, Nik, Hired_Date)
VALUES
    (uuid_generate_v4(), 'John Doe', 'Manager', '+6281234567890', '1234567890123456', '2023-01-01'),
    (uuid_generate_v4(), 'Jane Smith', 'Salesperson', '+6282345678901', '2345678901234567', '2023-02-01'),
    (uuid_generate_v4(), 'Robert Johnson', 'Warehouse Staff', '+6283456789012', '3456789012345678', '2023-03-01');

-- Insert Wages for Employees
INSERT INTO
    Wage (Id, Employee_Id, Total_Amount, Month, Year)
VALUES
    (uuid_generate_v4(), (SELECT Id FROM Employee WHERE Name = 'John Doe'), 5000000, 4, 2023),
    (uuid_generate_v4(), (SELECT Id FROM Employee WHERE Name = 'Jane Smith'), 4000000, 4, 2023),
    (uuid_generate_v4(), (SELECT Id FROM Employee WHERE Name = 'Robert Johnson'), 3500000, 4, 2023);

-- Insert Wage Details
INSERT INTO
    Wage_Detail (Id, Wage_Id, Component_Name, Description, Amount)
VALUES
    (uuid_generate_v4(), (SELECT Id FROM Wage WHERE Employee_Id = (SELECT Id FROM Employee WHERE Name = 'John Doe')), 'Basic Salary', 'Monthly base pay', 4000000),
    (uuid_generate_v4(), (SELECT Id FROM Wage WHERE Employee_Id = (SELECT Id FROM Employee WHERE Name = 'John Doe')), 'Bonus', 'Performance bonus', 1000000),
    (uuid_generate_v4(), (SELECT Id FROM Wage WHERE Employee_Id = (SELECT Id FROM Employee WHERE Name = 'Jane Smith')), 'Basic Salary', 'Monthly base pay', 3500000),
    (uuid_generate_v4(), (SELECT Id FROM Wage WHERE Employee_Id = (SELECT Id FROM Employee WHERE Name = 'Jane Smith')), 'Bonus', 'Sales target bonus', 500000),
    (uuid_generate_v4(), (SELECT Id FROM Wage WHERE Employee_Id = (SELECT Id FROM Employee WHERE Name = 'Robert Johnson')), 'Basic Salary', 'Monthly base pay', 3500000);

-- ================================================
-- Seed Data for Inventory
-- ================================================
-- Insert Categories
INSERT INTO
    Category (Id, Name, Description)
VALUES
    (uuid_generate_v4(), 'Electronics', 'Electronic devices and components'),
    (uuid_generate_v4(), 'Office Supplies', 'Items used in office settings'),
    (uuid_generate_v4(), 'Furniture', 'Office and home furniture');

-- Insert Units
INSERT INTO
    Unit (Id, Name, Description)
VALUES
    (uuid_generate_v4(), 'Piece', 'Individual item'),
    (uuid_generate_v4(), 'Box', 'Box containing multiple items'),
    (uuid_generate_v4(), 'Kilogram', 'Weight measurement');

-- Insert Products
INSERT INTO
    Product (Id, Name, Description, Category_Id, Unit_Id)
VALUES
    (uuid_generate_v4(), 'Laptop', 'Business laptop', (SELECT Id FROM Category WHERE Name = 'Electronics'), (SELECT Id FROM Unit WHERE Name = 'Piece')),
    (uuid_generate_v4(), 'Office Chair', 'Ergonomic office chair', (SELECT Id FROM Category WHERE Name = 'Furniture'), (SELECT Id FROM Unit WHERE Name = 'Piece')),
    (uuid_generate_v4(), 'Printer Paper', 'A4 printing paper', (SELECT Id FROM Category WHERE Name = 'Office Supplies'), (SELECT Id FROM Unit WHERE Name = 'Box'));

-- Insert Storages
INSERT INTO
    Storage (Id, Name, Location)
VALUES
    (uuid_generate_v4(), 'Main Warehouse', 'Jakarta'),
    (uuid_generate_v4(), 'Branch Storage', 'Bandung'),
    (uuid_generate_v4(), 'Retail Location', 'Jakarta Mall');

-- Insert Customers
INSERT INTO
    Customer (Id, Name, Address, Telephone)
VALUES
    (uuid_generate_v4(), 'PT ABC', 'Jakarta Business District', '+6280123456789'),
    (uuid_generate_v4(), 'CV XYZ', 'Bandung Industrial Area', '+6280234567890'),
    (uuid_generate_v4(), 'UD Maju Jaya', 'Surabaya Commercial Center', '+6280345678901');

-- Insert Suppliers
INSERT INTO
    Supplier (Id, Name, Address, Telephone)
VALUES
    (uuid_generate_v4(), 'Tech Supplies Inc.', 'Singapore Technology Park', '+6594567890123'),
    (uuid_generate_v4(), 'Office Depot', 'Jakarta Business Center', '+6281234987654'),
    (uuid_generate_v4(), 'Furniture World', 'Bandung Furniture District', '+6282345098765');

-- ================================================
-- Seed Data for Purchase
-- ================================================
-- Initialize document counters
INSERT INTO
    Document_Counter (Document_Type, Year, Month, Day, Counter)
VALUES
    ('PO', 2023, 4, 21, 0),
    ('SO', 2023, 4, 21, 0),
    ('SI', 2023, 4, 21, 0),
    ('DN', 2023, 4, 21, 0);


-- Insert Purchase Orders
DO $$
DECLARE
    laptop_id UUID;
    chair_id UUID;
    paper_id UUID;
    supplier_tech_id UUID;
    supplier_office_id UUID;
    purchase_order_id UUID;
    admin_id UUID;
    main_storage_id UUID;
BEGIN
    -- Get product IDs
    SELECT Id INTO laptop_id FROM Product WHERE Name = 'Laptop';
    SELECT Id INTO chair_id FROM Product WHERE Name = 'Office Chair';
    SELECT Id INTO paper_id FROM Product WHERE Name = 'Printer Paper';
    
    -- Get supplier IDs
    SELECT Id INTO supplier_tech_id FROM Supplier WHERE Name = 'Tech Supplies Inc.';
    SELECT Id INTO supplier_office_id FROM Supplier WHERE Name = 'Office Depot';
    
    -- Get user ID
    SELECT Id INTO admin_id FROM Appuser WHERE Username = 'admin';
    
    -- Get storage ID
    SELECT Id INTO main_storage_id FROM Storage WHERE Name = 'Main Warehouse';
    
    -- 1. Create a completed purchase order (cash payment)
    purchase_order_id := uuid_generate_v4();
    
    -- Insert purchase order
    INSERT INTO Purchase_Order (
        Id, Serial_Id, Supplier_Id, Order_Date, 
        Payment_Method, Status, Total_Amount, Created_By
    ) VALUES (
        purchase_order_id, 'PO-20230421-0001', supplier_tech_id, 
        '2023-04-21', 'cash', 'completed', 15000000, admin_id
    );
    
    -- Insert purchase order details
    INSERT INTO Purchase_Order_Detail (
        Id, Purchase_Order_Id, Product_Id, Requested_Quantity, Unit_Price
    ) VALUES
        (uuid_generate_v4(), purchase_order_id, laptop_id, 5, 3000000);
        
    -- Since this PO is completed, create product batches
    WITH batch_insert AS (
        INSERT INTO Product_Batch (
            Id, Sku, Product_Id, Purchase_Order_Id, Initial_Quantity, Current_Quantity, Unit_Price
        ) VALUES (
            uuid_generate_v4(), 'LAP-TEC210423-0001', laptop_id, purchase_order_id, 5, 5, 3000000
        ) RETURNING Id
    )
    INSERT INTO Batch_Storage (
        Id, Batch_Id, Storage_Id, Quantity
    ) SELECT 
        uuid_generate_v4(), Id, main_storage_id, 5
    FROM batch_insert;
    
    -- 2. Create an ordered purchase order (credit payment)
    purchase_order_id := uuid_generate_v4();
    
    -- Insert purchase order with payment due date
    INSERT INTO Purchase_Order (
        Id, Serial_Id, Supplier_Id, Order_Date, 
        Payment_Method, Payment_Due_Date, Status, Total_Amount, Created_By
    ) VALUES (
        purchase_order_id, 'PO-20230421-0002', supplier_office_id, 
        '2023-04-21', 'credit', '2023-05-21', 'ordered', 5500000, admin_id
    );
    
    -- Insert purchase order details
    INSERT INTO Purchase_Order_Detail (
        Id, Purchase_Order_Id, Product_Id, Requested_Quantity, Unit_Price
    ) VALUES
        (uuid_generate_v4(), purchase_order_id, chair_id, 10, 500000),
        (uuid_generate_v4(), purchase_order_id, paper_id, 20, 25000);
END $$;

-- ================================================
-- Seed Data for Sales
-- ================================================
-- Insert Sales Orders
DO $$
DECLARE
    laptop_id UUID;
    customer_abc_id UUID;
    customer_xyz_id UUID;
    admin_id UUID;
    sales_order_id UUID;
    sales_invoice_id UUID;
    batch_storage_id UUID;
BEGIN
    -- Get product and batch storage IDs
    SELECT Product_Id INTO laptop_id FROM Product_Batch WHERE SKU = 'LAP-TEC210423-0001';
    SELECT bs.Id INTO batch_storage_id 
    FROM Batch_Storage bs 
    JOIN Product_Batch pb ON bs.Batch_Id = pb.Id 
    WHERE pb.SKU = 'LAP-TEC210423-0001';
    
    -- Get customer IDs
    SELECT Id INTO customer_abc_id FROM Customer WHERE Name = 'PT ABC';
    SELECT Id INTO customer_xyz_id FROM Customer WHERE Name = 'CV XYZ';
    
    -- Get user ID
    SELECT Id INTO admin_id FROM Appuser WHERE Username = 'admin';
    
    -- 1. Create a sales order that's already invoiced (cash payment)
    sales_order_id := uuid_generate_v4();
    
    -- Insert sales order
    INSERT INTO Sales_Order (
        Id, Serial_Id, Customer_Id, Order_Date, 
        Payment_Method, Status, Total_Amount, Created_By
    ) VALUES (
        sales_order_id, 'SO-20230421-0001', customer_abc_id, 
        '2023-04-21', 'cash', 'invoice', 6600000, admin_id
    );
    
    -- Insert sales order details
    INSERT INTO Sales_Order_Detail (
        Id, Sales_Order_Id, Product_Id, Batch_Storage_Id, Quantity, Unit_Price
    ) VALUES
        (uuid_generate_v4(), sales_order_id, laptop_id, batch_storage_id, 2, 3300000);
        
    -- Create a sales invoice
    sales_invoice_id := uuid_generate_v4();
    INSERT INTO Sales_Invoice (
        Id, Serial_Id, Sales_Order_Id, Invoice_Date, Total_Amount, Created_By
    ) VALUES (
        sales_invoice_id, 'SI-20230421-0001', sales_order_id, '2023-04-21', 6600000, admin_id
    );
    
    -- Update batch storage quantity (subtract 2)
    UPDATE Batch_Storage 
    SET Quantity = Quantity - 2
    WHERE Id = batch_storage_id;
    
    UPDATE Product_Batch
    SET Current_Quantity = Current_Quantity - 2
    WHERE Id = (SELECT Batch_Id FROM Batch_Storage WHERE Id = batch_storage_id);
    
    -- 2. Create a sales order with "order" status (paylater payment)
    sales_order_id := uuid_generate_v4();
    
    -- Insert sales order with payment due date
    INSERT INTO Sales_Order (
        Id, Serial_Id, Customer_Id, Order_Date, 
        Payment_Method, Payment_Due_Date, Status, Total_Amount, Created_By
    ) VALUES (
        sales_order_id, 'SO-20230421-0002', customer_xyz_id, 
        '2023-04-21', 'paylater', '2023-05-21', 'order', 3300000, admin_id
    );
    
    -- Insert sales order details
    INSERT INTO Sales_Order_Detail (
        Id, Sales_Order_Id, Product_Id, Batch_Storage_Id, Quantity, Unit_Price
    ) VALUES
        (uuid_generate_v4(), sales_order_id, laptop_id, batch_storage_id, 1, 3300000);
        
    -- Update batch storage quantity (subtract 1)
    UPDATE Batch_Storage 
    SET Quantity = Quantity - 1
    WHERE Id = batch_storage_id;
    
    UPDATE Product_Batch
    SET Current_Quantity = Current_Quantity - 1
    WHERE Id = (SELECT Batch_Id FROM Batch_Storage WHERE Id = batch_storage_id);
END $$;