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
-- INSERT INTO
--     Document_Counter (Document_Type, Year, Month, Day, Counter)
-- VALUES
--     ('PO', 2023, 4, 21, 0),
--     ('SO', 2023, 4, 21, 0),
--     ('SI', 2023, 4, 21, 0),
--     ('DN', 2023, 4, 21, 0);


-- -- Insert Purchase Orders
-- DO $$
-- DECLARE
--     laptop_id UUID;
--     chair_id UUID;
--     paper_id UUID;
--     supplier_tech_id UUID;
--     supplier_office_id UUID;
--     supplier_furniture_id UUID;
--     purchase_order_id UUID;
--     purchase_order_detail_id UUID;
--     admin_id UUID;
--     main_storage_id UUID;
--     branch_storage_id UUID;
--     batch_id UUID;
-- BEGIN
--     -- Get product IDs
--     SELECT Id INTO laptop_id FROM Product WHERE Name = 'Laptop';
--     SELECT Id INTO chair_id FROM Product WHERE Name = 'Office Chair';
--     SELECT Id INTO paper_id FROM Product WHERE Name = 'Printer Paper';
    
--     -- Get supplier IDs
--     SELECT Id INTO supplier_tech_id FROM Supplier WHERE Name = 'Tech Supplies Inc.';
--     SELECT Id INTO supplier_office_id FROM Supplier WHERE Name = 'Office Depot';
--     SELECT Id INTO supplier_furniture_id FROM Supplier WHERE Name = 'Furniture World';
    
--     -- Get user ID
--     SELECT Id INTO admin_id FROM Appuser WHERE Username = 'admin';
    
--     -- Get storage IDs
--     SELECT Id INTO main_storage_id FROM Storage WHERE Name = 'Main Warehouse';
--     SELECT Id INTO branch_storage_id FROM Storage WHERE Name = 'Branch Storage';
    
--     -- 1. Create a completed purchase order (cash payment)
--     purchase_order_id := uuid_generate_v4();
    
--     -- Insert purchase order
--     INSERT INTO Purchase_Order (
--         Id, Serial_Id, Supplier_Id, Order_Date, 
--         Payment_Method, Status, Total_Amount, Created_By, Checked_By
--     ) VALUES (
--         purchase_order_id, 'PO-20230421-0001', supplier_tech_id, 
--         '2023-04-21', 'cash', 'completed', 15000000, admin_id, admin_id
--     );
    
--     -- Insert purchase order details
--     INSERT INTO Purchase_Order_Detail (
--         Id, Purchase_Order_Id, Product_Id, Requested_Quantity, Unit_Price
--     ) VALUES
--         (uuid_generate_v4(), purchase_order_id, laptop_id, 5, 3000000);
        
--     -- Since this PO is completed, create product batches
--     WITH batch_insert AS (
--         INSERT INTO Product_Batch (
--             Id, Sku, Product_Id, Purchase_Order_Id, Initial_Quantity, Current_Quantity, Unit_Price
--         ) VALUES (
--             uuid_generate_v4(), 'LAP-TEC210423-0001', laptop_id, purchase_order_id, 5, 5, 3000000
--         ) RETURNING Id
--     )
--     INSERT INTO Batch_Storage (
--         Id, Batch_Id, Storage_Id, Quantity
--     ) SELECT 
--         uuid_generate_v4(), Id, main_storage_id, 5
--     FROM batch_insert;
    
--     -- Record inventory log for the purchase
--     INSERT INTO Inventory_Log (
--         Batch_Id, Storage_Id, User_Id, Purchase_Order_Id, Action, Quantity, Description
--     )
--     SELECT 
--         pb.Id, bs.Storage_Id, admin_id, purchase_order_id, 'add', 5, 
--         'Pembelian Barang PO-20230421-0001'
--     FROM 
--         Product_Batch pb
--         JOIN Batch_Storage bs ON pb.Id = bs.Batch_Id
--     WHERE 
--         pb.Purchase_Order_Id = purchase_order_id;
        
--     -- Record financial transaction log for the purchase
--     INSERT INTO Financial_Transaction_Log (
--         User_Id, Amount, Type, Purchase_Order_Id, Description, Transaction_Date, Is_System
--     ) VALUES (
--         admin_id, 15000000, 'purchase', purchase_order_id, 
--         'Pembelian Barang PO-20230421-0001', '2023-04-21', true
--     );
    
--     -- 2. Create an ordered purchase order (credit payment)
--     purchase_order_id := uuid_generate_v4();
    
--     -- Insert purchase order with payment due date
--     INSERT INTO Purchase_Order (
--         Id, Serial_Id, Supplier_Id, Order_Date, 
--         Payment_Method, Payment_Due_Date, Status, Total_Amount, Created_By
--     ) VALUES (
--         purchase_order_id, 'PO-20230421-0002', supplier_office_id, 
--         '2023-04-21', 'credit', '2023-05-21', 'ordered', 5500000, admin_id
--     );
    
--     -- Insert purchase order details
--     INSERT INTO Purchase_Order_Detail (
--         Id, Purchase_Order_Id, Product_Id, Requested_Quantity, Unit_Price
--     ) VALUES
--         (uuid_generate_v4(), purchase_order_id, chair_id, 10, 500000),
--         (uuid_generate_v4(), purchase_order_id, paper_id, 20, 25000);
        
--     -- 3. Create a partially returned purchase order
--     purchase_order_id := uuid_generate_v4();
    
--     -- Insert purchase order
--     INSERT INTO Purchase_Order (
--         Id, Serial_Id, Supplier_Id, Order_Date, 
--         Payment_Method, Status, Total_Amount, Created_By, Checked_By
--     ) VALUES (
--         purchase_order_id, 'PO-20230421-0003', supplier_furniture_id, 
--         '2023-04-21', 'cash', 'partially_returned', 5000000, admin_id, admin_id
--     );
    
--     -- Insert purchase order detail
--     purchase_order_detail_id := uuid_generate_v4();
--     INSERT INTO Purchase_Order_Detail (
--         Id, Purchase_Order_Id, Product_Id, Requested_Quantity, Unit_Price
--     ) VALUES
--         (purchase_order_detail_id, purchase_order_id, chair_id, 10, 500000);
    
--     -- Create product batch for the order
--     INSERT INTO Product_Batch (
--         Id, Sku, Product_Id, Purchase_Order_Id, Initial_Quantity, Current_Quantity, Unit_Price
--     ) VALUES (
--         uuid_generate_v4(), 'CHR-FUR210423-0001', chair_id, purchase_order_id, 10, 8, 500000
--     ) RETURNING Id INTO batch_id;
    
--     -- Add to batch storage
--     INSERT INTO Batch_Storage (
--         Id, Batch_Id, Storage_Id, Quantity
--     ) VALUES (
--         uuid_generate_v4(), batch_id, branch_storage_id, 8
--     );
    
--     -- Create return record
--     INSERT INTO Purchase_Order_Return (
--         Id, Purchase_Order_Id, Product_Detail_Id, Return_Quantity, 
--         Reason, Status, Returned_By, Returned_At
--     ) VALUES (
--         uuid_generate_v4(), purchase_order_id, purchase_order_detail_id, 2,
--         'Damaged on arrival', 'returned', admin_id, '2023-04-22'
--     );
    
--     -- Record inventory log for the purchase and return
--     INSERT INTO Inventory_Log (
--         Batch_Id, Storage_Id, User_Id, Purchase_Order_Id, Action, Quantity, Description
--     ) VALUES
--     (
--         batch_id, branch_storage_id, admin_id, purchase_order_id, 'add', 10, 
--         'Pembelian Barang PO-20230421-0003'
--     ),
--     (
--         batch_id, branch_storage_id, admin_id, purchase_order_id, 'return', 2, 
--         'Retur Pembelian PO-20230421-0003'
--     );
    
--     -- Record financial transactions for purchase and return
--     INSERT INTO Financial_Transaction_Log (
--         User_Id, Amount, Type, Purchase_Order_Id, Description, Transaction_Date, Is_System
--     ) VALUES
--     (
--         admin_id, 5000000, 'purchase', purchase_order_id, 
--         'Pembelian Barang PO-20230421-0003', '2023-04-21', true
--     ),
--     (
--         admin_id, -1000000, 'purchase_return', purchase_order_id, 
--         'Retur Pembelian PO-20230421-0003', '2023-04-22', true
--     );
-- END $$;

-- -- ================================================
-- -- Seed Data for Sales
-- -- ================================================
-- -- Insert Sales Orders
-- DO $$
-- DECLARE
--     laptop_id UUID;
--     chair_id UUID;
--     customer_abc_id UUID;
--     customer_xyz_id UUID;
--     customer_ud_id UUID;
--     admin_id UUID;
--     sales_order_id UUID;
--     sales_invoice_id UUID;
--     delivery_note_id UUID;
--     laptop_batch_storage_id UUID;
--     chair_batch_storage_id UUID;
--     sales_detail_id UUID;
-- BEGIN
--     -- Get product IDs
--     SELECT Product_Id INTO laptop_id FROM Product_Batch WHERE SKU = 'LAP-TEC210423-0001';
--     SELECT Product_Id INTO chair_id FROM Product_Batch WHERE SKU = 'CHR-FUR210423-0001';
    
--     -- Get batch storage IDs
--     SELECT bs.Id INTO laptop_batch_storage_id 
--     FROM Batch_Storage bs 
--     JOIN Product_Batch pb ON bs.Batch_Id = pb.Id 
--     WHERE pb.SKU = 'LAP-TEC210423-0001';
    
--     SELECT bs.Id INTO chair_batch_storage_id 
--     FROM Batch_Storage bs 
--     JOIN Product_Batch pb ON bs.Batch_Id = pb.Id 
--     WHERE pb.SKU = 'CHR-FUR210423-0001';
    
--     -- Get customer IDs
--     SELECT Id INTO customer_abc_id FROM Customer WHERE Name = 'PT ABC';
--     SELECT Id INTO customer_xyz_id FROM Customer WHERE Name = 'CV XYZ';
--     SELECT Id INTO customer_ud_id FROM Customer WHERE Name = 'UD Maju Jaya';
    
--     -- Get user ID
--     SELECT Id INTO admin_id FROM Appuser WHERE Username = 'admin';
    
--     -- 1. Create a complete sales order with invoice and delivery note (cash payment)
--     sales_order_id := uuid_generate_v4();
    
--     -- Insert sales order
--     INSERT INTO Sales_Order (
--         Id, Serial_Id, Customer_Id, Order_Date, 
--         Payment_Method, Status, Total_Amount, Created_By
--     ) VALUES (
--         sales_order_id, 'SO-20230421-0001', customer_abc_id, 
--         '2023-04-21', 'cash', 'invoice', 6600000, admin_id
--     );
    
--     -- Insert sales order details
--     INSERT INTO Sales_Order_Detail (
--         Id, Sales_Order_Id, Batch_Storage_Id, Quantity, Unit_Price
--     ) VALUES
--         (uuid_generate_v4(), sales_order_id, laptop_batch_storage_id, 2, 3300000);
        
--     -- Create a sales invoice
--     sales_invoice_id := uuid_generate_v4();
--     INSERT INTO Sales_Invoice (
--         Id, Serial_Id, Sales_Order_Id, Invoice_Date, Total_Amount, Created_By
--     ) VALUES (
--         sales_invoice_id, 'SI-20230421-0001', sales_order_id, '2023-04-21', 6600000, admin_id
--     );
    
--     -- Create a delivery note
--     delivery_note_id := uuid_generate_v4();
--     INSERT INTO Delivery_Note (
--         Id, Serial_Id, Sales_Order_Id, Sales_Invoice_Id, 
--         Delivery_Date, Driver_Name, Recipient_Name, Created_By
--     ) VALUES (
--         delivery_note_id, 'DN-20230421-0001', sales_order_id, sales_invoice_id,
--         '2023-04-22', 'Budi Santoso', 'John PT ABC', admin_id
--     );
    
--     -- Update batch storage quantity (subtract 2)
--     UPDATE Batch_Storage 
--     SET Quantity = Quantity - 2
--     WHERE Id = laptop_batch_storage_id;
    
--     UPDATE Product_Batch
--     SET Current_Quantity = Current_Quantity - 2
--     WHERE Id = (SELECT Batch_Id FROM Batch_Storage WHERE Id = laptop_batch_storage_id);
    
--     -- Record inventory log for the sales transaction
--     INSERT INTO Inventory_Log (
--         Batch_Id, Storage_Id, User_Id, Sales_Order_Id, Action, Quantity, Description
--     ) VALUES (
--         (SELECT Batch_Id FROM Batch_Storage WHERE Id = laptop_batch_storage_id),
--         (SELECT Storage_Id FROM Batch_Storage WHERE Id = laptop_batch_storage_id),
--         admin_id, sales_order_id, 'remove', 2, 'Pembuatan faktur SI-20230421-0001'
--     );
    
--     -- Record financial transaction log
--     INSERT INTO Financial_Transaction_Log (
--         User_Id, Amount, Type, Sales_Order_Id, Description, Transaction_Date, Is_System
--     ) VALUES (
--         admin_id, 6600000, 'sales_invoice', sales_order_id, 
--         'Pembuatan faktur penjualan SI-20230421-0001', '2023-04-21', true
--     );
    
--     -- 2. Create a sales order with "order" status (paylater payment)
--     sales_order_id := uuid_generate_v4();
    
--     -- Insert sales order with payment due date
--     INSERT INTO Sales_Order (
--         Id, Serial_Id, Customer_Id, Order_Date, 
--         Payment_Method, Payment_Due_Date, Status, Total_Amount, Created_By
--     ) VALUES (
--         sales_order_id, 'SO-20230421-0002', customer_xyz_id, 
--         '2023-04-21', 'paylater', '2023-05-21', 'order', 3300000, admin_id
--     );
    
--     -- Insert sales order details
--     sales_detail_id := uuid_generate_v4();
--     INSERT INTO Sales_Order_Detail (
--         Id, Sales_Order_Id, Batch_Storage_Id, Quantity, Unit_Price
--     ) VALUES
--         (sales_detail_id, sales_order_id, laptop_batch_storage_id, 1, 3300000);
        
--     -- Update batch storage quantity (subtract 1)
--     UPDATE Batch_Storage 
--     SET Quantity = Quantity - 1
--     WHERE Id = laptop_batch_storage_id;
    
--     UPDATE Product_Batch
--     SET Current_Quantity = Current_Quantity - 1
--     WHERE Id = (SELECT Batch_Id FROM Batch_Storage WHERE Id = laptop_batch_storage_id);
    
--     -- 3. Create a sales order with chairs that has return
--     sales_order_id := uuid_generate_v4();
    
--     -- Insert sales order
--     INSERT INTO Sales_Order (
--         Id, Serial_Id, Customer_Id, Order_Date, 
--         Payment_Method, Status, Total_Amount, Created_By
--     ) VALUES (
--         sales_order_id, 'SO-20230421-0003', customer_ud_id, 
--         '2023-04-21', 'cash', 'partially_returned', 3000000, admin_id
--     );
    
--     -- Insert sales order details
--     sales_detail_id := uuid_generate_v4();
--     INSERT INTO Sales_Order_Detail (
--         Id, Sales_Order_Id, Batch_Storage_Id, Quantity, Unit_Price
--     ) VALUES
--         (sales_detail_id, sales_order_id, chair_batch_storage_id, 5, 600000);
    
--     -- Create a sales invoice
--     sales_invoice_id := uuid_generate_v4();
--     INSERT INTO Sales_Invoice (
--         Id, Serial_Id, Sales_Order_Id, Invoice_Date, Total_Amount, Created_By
--     ) VALUES (
--         sales_invoice_id, 'SI-20230421-0002', sales_order_id, '2023-04-21', 3000000, admin_id
--     );
    
--     -- Update batch storage quantity (subtract 5)
--     UPDATE Batch_Storage 
--     SET Quantity = Quantity - 5
--     WHERE Id = chair_batch_storage_id;
    
--     UPDATE Product_Batch
--     SET Current_Quantity = Current_Quantity - 5
--     WHERE Id = (SELECT Batch_Id FROM Batch_Storage WHERE Id = chair_batch_storage_id);
    
--     -- Create return record
--     INSERT INTO Sales_Order_Return (
--         Id, Return_Source, Sales_Order_Id, Sales_Detail_Id, 
--         Return_Quantity, Remaining_Quantity, Return_Reason, 
--         Return_Status, Returned_By, Returned_At
--     ) VALUES (
--         uuid_generate_v4(), 'invoice', sales_order_id, sales_detail_id, 
--         2, 3, 'Wrong size ordered', 'completed', admin_id, '2023-04-23'
--     );
    
--     -- Update batch storage quantity for the return (add 2 back)
--     UPDATE Batch_Storage 
--     SET Quantity = Quantity + 2
--     WHERE Id = chair_batch_storage_id;
    
--     UPDATE Product_Batch
--     SET Current_Quantity = Current_Quantity + 2
--     WHERE Id = (SELECT Batch_Id FROM Batch_Storage WHERE Id = chair_batch_storage_id);
    
--     -- Record inventory logs
--     INSERT INTO Inventory_Log (
--         Batch_Id, Storage_Id, User_Id, Sales_Order_Id, Action, Quantity, Description
--     ) VALUES
--     (
--         (SELECT Batch_Id FROM Batch_Storage WHERE Id = chair_batch_storage_id),
--         (SELECT Storage_Id FROM Batch_Storage WHERE Id = chair_batch_storage_id),
--         admin_id, sales_order_id, 'remove', 5, 'Pembuatan faktur SI-20230421-0002'
--     ),
--     (
--         (SELECT Batch_Id FROM Batch_Storage WHERE Id = chair_batch_storage_id),
--         (SELECT Storage_Id FROM Batch_Storage WHERE Id = chair_batch_storage_id),
--         admin_id, sales_order_id, 'return', 2, 'Retur penjualan SO-20230421-0003'
--     );
    
--     -- Record financial transactions
--     INSERT INTO Financial_Transaction_Log (
--         User_Id, Amount, Type, Sales_Order_Id, Description, Transaction_Date, Is_System
--     ) VALUES
--     (
--         admin_id, 3000000, 'sales_invoice', sales_order_id, 
--         'Pembuatan faktur penjualan SI-20230421-0002', '2023-04-21', true
--     ),
--     (
--         admin_id, -1200000, 'sales_return', sales_order_id, 
--         'Retur penjualan SO-20230421-0003', '2023-04-23', true
--     );
-- END $$;
