CREATE DATABASE IF NOT EXISTS liga_db;
USE liga_db;

CREATE TABLE IF NOT EXISTS matches (
  id INT AUTO_INCREMENT PRIMARY KEY,
  home_team VARCHAR(100),
  away_team VARCHAR(100),
  match_date DATE
);

INSERT INTO matches (home_team, away_team, match_date) VALUES
('Real Madrid', 'Barcelona', '2025-04-01'),
('Atletico', 'Sevilla', '2025-04-02'),
('Valencia', 'Villarreal', '2025-04-03'),
('Betis', 'Celta de Vigo', '2025-04-04'),
('Osasuna', 'Getafe', '2025-04-05');