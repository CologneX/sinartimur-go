-- Data seed untuk schema Inventory dengan UUID valid

--------------------------------------------
-- Seed Data untuk Tabel Category
--------------------------------------------
INSERT INTO Category (Id, Name, Description)
VALUES ('a1b2c3d4-5678-90ab-cdef-1234567890ab', 'Elektronik',
        'Kategori untuk produk-produk elektronik seperti TV, smartphone, dan komputer'),
       ('b2c3d4e5-6789-0abc-def1-234567890abc', 'Perabotan',
        'Kategori untuk perabotan rumah dan kantor seperti meja, kursi, dan lemari'),
       ('c3d4e5f6-7890-1bcd-ef12-34567890abcd', 'Alat Tulis',
        'Kategori untuk keperluan sekolah dan kantor seperti pulpen, pensil, dan buku');

--------------------------------------------
-- Seed Data untuk Tabel Unit
--------------------------------------------
INSERT INTO Unit (Id, Name, Description)
VALUES ('d4e5f6a7-8901-2cde-f123-4567890abcde', 'PCS', 'Satuan per buah'),
       ('e5f6a7b8-9012-3def-1234-567890abcdef', 'KG', 'Satuan kilogram');

--------------------------------------------
-- Seed Data untuk Tabel Storages
--------------------------------------------
INSERT INTO Storage (Id, Name, Location)
VALUES ('f6a7b8c9-0123-4def-2345-67890abcdef1', 'Gudang Utama', 'Lokasi: Jalan Merdeka No.1, Jakarta'),
       ('a7b8c9d0-1234-5ef0-3456-7890abcdef12', 'Gudang Cabang A', 'Lokasi: Jalan Sudirman No.10, Bandung');

--------------------------------------------
-- Seed Data untuk Tabel Products
--------------------------------------------
INSERT INTO Product (Id, Name, Description, Price, Category_Id, Unit_Id)
VALUES ('b8c9d0e1-2345-6f01-4567-890abcdef123', 'Televisi LED 32 inch', 'Televisi LED 32 inch dengan resolusi tinggi',
        2500000, 'a1b2c3d4-5678-90ab-cdef-1234567890ab', 'd4e5f6a7-8901-2cde-f123-4567890abcde'),
       ('c9d0e1f2-3456-7f12-5678-90abcdef1234', 'Kursi Kantor Ergonomis', 'Kursi kantor dengan desain ergonomis',
        750000, 'b2c3d4e5-6789-0abc-def1-234567890abc', 'd4e5f6a7-8901-2cde-f123-4567890abcde'),
       ('d0e1f203-4567-8f23-6789-0abcdef12345', 'Pulpen Gel', 'Pulpen gel dengan tinta tahan lama', 5000,
        'c3d4e5f6-7890-1bcd-ef12-34567890abcd', 'd4e5f6a7-8901-2cde-f123-4567890abcde');

--------------------------------------------
-- Seed Data untuk Tabel Inventory
--------------------------------------------
INSERT INTO Inventory (Id, Product_Id, Storage_Id, Quantity, Minimum_Quantity)
VALUES ('e1f20304-5678-9f34-7890-abcdef123456', 'b8c9d0e1-2345-6f01-4567-890abcdef123',
        'f6a7b8c9-0123-4def-2345-67890abcdef1', 10, 5),
       ('f2030456-6789-af45-8901-bcdef1234567', 'c9d0e1f2-3456-7f12-5678-90abcdef1234',
        'a7b8c9d0-1234-5ef0-3456-7890abcdef12', 20, 10),
       ('10304567-789a-bf56-9012-cdef12345678', 'd0e1f203-4567-8f23-6789-0abcdef12345',
        'f6a7b8c9-0123-4def-2345-67890abcdef1', 100, 50);

--------------------------------------------
-- Seed Data untuk Tabel Inventory_Logs
--------------------------------------------
INSERT INTO Inventory_Log (Id, Inventory_Id, User_Id, Action, Quantity, Description)
VALUES ('20345678-89ab-cf67-0123-def123456789', 'e1f20304-5678-9f34-7890-abcdef123456', NULL, 'add', 5,
        'Penambahan stok awal untuk Televisi LED 32 inch'),
       ('30456789-9abc-df78-1234-ef1234567890', 'f2030456-6789-af45-8901-bcdef1234567', NULL, 'remove', 2,
        'Pengurangan stok karena kerusakan pada Kursi Kantor Ergonomis');
