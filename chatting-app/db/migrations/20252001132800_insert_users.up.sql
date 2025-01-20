INSERT INTO users (username, password, email, full_name, phone_number, is_active, role) VALUES
('johndoe', 'johndoe', 'johndoe@example.com', 'John Doe', '+1234567890', TRUE, 'user'),
('janedoe', 'janedoe', 'janedoe@example.com', 'Jane Doe', NULL, TRUE, 'user'),
('admin', 'admin', 'admin@example.com', 'Admin User', '+9876543210', TRUE, 'admin'),
('guest', 'guest', NULL, 'Guest User', NULL, FALSE, 'guest'),
('alice', 'alice', 'alice@example.com', 'Alice Johnson', '+1122334455', TRUE, 'user'),
('bob', 'bob', 'bob@example.com', 'Bob Smith', '+9988776655', TRUE, 'user'),
('charlie', 'charlie', NULL, 'Charlie Brown', NULL, FALSE, 'guest'),
('david', 'david', 'david@example.com', 'David Lee', '+6677889900', TRUE, 'user'),
('eve', 'eve', 'eve@example.com', 'Eve Adams', '+5566778899', TRUE, 'user'),
('mallory', 'mallory', 'mallory@example.com', 'Mallory Black', '+4455667788', FALSE, 'admin');