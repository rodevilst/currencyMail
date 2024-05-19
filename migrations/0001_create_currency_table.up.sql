CREATE TABLE IF NOT EXISTS currency (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ccy VARCHAR(10) NOT NULL,
    base_ccy VARCHAR(10) NOT NULL,
    buy VARCHAR(20) NOT NULL,
    sale VARCHAR(20) NOT NULL,
    UNIQUE KEY unique_currency (ccy, base_ccy)
);
