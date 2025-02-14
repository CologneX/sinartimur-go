package inventory

//
//
//Create Table Storage
//(
//Id         Uuid Primary Key Default Uuid_Generate_V4(),
//Name       VARCHAR(255) Not Null,
//Location   TEXT         Not Null,
//Created_At Timestamptz      Default Current_Timestamp,
//Updated_At Timestamptz      Default Current_Timestamp,
//Deleted_At Timestamptz      Default Null
//);
//
//Create Table Inventory
//(
//Id               Uuid Primary Key Default Uuid_Generate_V4(),
//Product_Id       Uuid References Product (Id) On Delete Cascade,
//Storage_Id       Uuid References Storage (Id) On Delete Cascade,
//Quantity         INT Not Null     Default 0,
//Minimum_Quantity INT Not Null     Default 0,
//Created_At       Timestamptz      Default Current_Timestamp,
//Updated_At       Timestamptz      Default Current_Timestamp
//);
//
//
//Create Table Inventory_Log
//(
//Id           Uuid Primary Key Default Uuid_Generate_V4(),
//Inventory_Id Uuid References Inventory (Id) On Delete Cascade,
//User_Id      Uuid References AppUser (Id),
//Action       VARCHAR(50) Not Null, -- e.g., "add", "remove", "transfer"
//Quantity     INT         Not Null,
//Log_Date     Timestamptz      Default Current_Timestamp,
//Description  TEXT,
//Created_At   Timestamptz      Default Current_Timestamp,
//Updated_At   Timestamptz      Default Current_Timestamp
//);

type Storage struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type Inventory struct {
	ID              string `json:"id"`
	ProductID       string `json:"product_id"`
	StorageID       string `json:"storage_id"`
	Quantity        int    `json:"quantity"`
	MinimumQuantity int    `json:"minimum_quantity"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type InventoryLog struct {
	ID          string `json:"id"`
	InventoryID string `json:"inventory_id"`
	UserID      string `json:"user_id"`
	Action      string `json:"action"`
	Quantity    int    `json:"quantity"`
	LogDate     string `json:"log_date"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
