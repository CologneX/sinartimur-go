-- Enable UUID Extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table: Admin
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       username VARCHAR(100) UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       is_active BOOLEAN DEFAULT TRUE, -- Status user aktif/nonaktif
                       created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE roles (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(50) UNIQUE NOT NULL,
                       description TEXT,
                       created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_roles (
                            id SERIAL PRIMARY KEY,
                            user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                            role_id INT REFERENCES roles(id) ON DELETE CASCADE,
                            assigned_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                            UNIQUE (user_id, role_id)
);

-- Table: HR
CREATE TABLE employees (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           name VARCHAR(150) NOT NULL,
                           position VARCHAR(100) NOT NULL, -- Posisi karyawan, e.g., "Manager", "Staff"
                           hired_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                           deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE employees (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           name VARCHAR(150) NOT NULL,
                           position VARCHAR(100) NOT NULL, -- Posisi karyawan, e.g., "Manager", "Staff"
                           hired_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                           deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE wages (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       employee_id UUID REFERENCES employees(id) ON DELETE CASCADE,
                       total_amount NUMERIC(12, 2) NOT NULL,
                       period_start DATE NOT NULL,
                       period_end DATE NOT NULL,
                       created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE wage_details (
                              id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                              wage_id UUID REFERENCES wages(id) ON DELETE CASCADE,
                              component_name VARCHAR(100) NOT NULL,
                              description TEXT,
                              amount NUMERIC(12, 2) NOT NULL,
                              created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                              updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                              deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Table: Financial Transactions
CREATE TABLE financial_transactions (
                                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                        user_id UUID REFERENCES users(id),
                                        amount NUMERIC(15, 2) NOT NULL,
                                        type VARCHAR(50) NOT NULL,
                                        description TEXT,
                                        transaction_date TIMESTAMPTZ NOT NULL,
                                        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                                        edited_at TIMESTAMPTZ DEFAULT NULL,
                                        deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Table: Inventory
CREATE TABLE inventory (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           name VARCHAR(255) NOT NULL,
                           description TEXT,
                           quantity INT NOT NULL DEFAULT 0,
                           minimum_quantity INT NOT NULL DEFAULT 0,
                           created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                           edited_at TIMESTAMPTZ DEFAULT NULL,
                           deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Table: Inventory Logs (for tracking stock movements)
CREATE TABLE inventory_logs (
                                id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                inventory_id UUID REFERENCES inventory(id),
                                user_id UUID REFERENCES users(id),
                                action VARCHAR(50) NOT NULL, -- e.g., add, remove
                                quantity INT NOT NULL,
                                log_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                                edited_at TIMESTAMPTZ DEFAULT NULL,
                                deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Table: Orders
CREATE TABLE orders (
                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                        customer_name VARCHAR(255) NOT NULL,
                        order_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                        status VARCHAR(50) NOT NULL,
                        total_amount NUMERIC(15, 2) NOT NULL,
                        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                        edited_at TIMESTAMPTZ DEFAULT NULL,
                        deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Table: Order Items
CREATE TABLE order_items (
                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                             order_id UUID REFERENCES orders(id),
                             inventory_id UUID REFERENCES inventory(id),
                             quantity INT NOT NULL,
                             price NUMERIC(15, 2) NOT NULL,
                             created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                             edited_at TIMESTAMPTZ DEFAULT NULL,
                             deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Indexes to improve query performance
CREATE INDEX idx_financial_transactions_user_id ON financial_transactions(user_id);
CREATE INDEX idx_inventory_name ON inventory(name);
CREATE INDEX idx_inventory_logs_inventory_id ON inventory_logs(inventory_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_employees_name ON employees(name);
CREATE INDEX idx_employees_position ON employees(position);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX idx_employees_name ON employees(name);
CREATE INDEX idx_employees_position ON employees(position);
CREATE INDEX idx_wages_employee_id ON wages(employee_id);
CREATE INDEX idx_wages_period ON wages(period_start, period_end);
CREATE INDEX idx_wage_details_wage_id ON wage_details(wage_id);
CREATE INDEX idx_wage_details_component_name ON wage_details(component_name);
