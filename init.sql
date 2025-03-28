CREATE DATABASE IF NOT EXISTS liga_db;
USE liga_db;

DROP TABLE IF EXISTS matches;

CREATE TABLE matches (
  id INT AUTO_INCREMENT PRIMARY KEY,
  home_team VARCHAR(100),
  away_team VARCHAR(100),
  match_date DATE,
  goals INT DEFAULT 0,
  yellow_cards INT DEFAULT 0,
  red_cards INT DEFAULT 0,
  extra_time VARCHAR(20) DEFAULT ''
);

INSERT INTO matches (home_team, away_team, match_date) VALUES
('Real Madrid', 'Barcelona', '2025-04-01'),
('Atletico', 'Sevilla', '2025-04-02'),
('Valencia', 'Villarreal', '2025-04-03'),
('Betis', 'Celta de Vigo', '2025-04-04'),
('Osasuna', 'Getafe', '2025-04-05'),
('Girona', 'Mallorca', '2025-04-06'),
('Alavés', 'Granada', '2025-04-07'),
('Espanyol', 'Levante', '2025-04-08'),
('Rayo Vallecano', 'Zaragoza', '2025-04-09'),
('Real Sociedad', 'Elche', '2025-04-10'),
('Las Palmas', 'Cádiz', '2025-04-11'),
('Eibar', 'Tenerife', '2025-04-12'),
('Sporting', 'Oviedo', '2025-04-13'),
('Huesca', 'Leganés', '2025-04-14'),
('Mirandés', 'Alcorcón', '2025-04-15');