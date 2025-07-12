CREATE TABLE IF NOT EXISTS event_detail (
  id INT PRIMARY KEY AUTO_INCREMENT,
  title varchar(255) NOT NULL,
  organizer_id INT NOT NULL COMMENT 'user id of the organizer',
  duration_minutes INT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);