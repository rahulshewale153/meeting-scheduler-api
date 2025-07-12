CREATE TABLE IF NOT EXISTS user_availability (
  id INT PRIMARY KEY AUTO_INCREMENT,
  event_id INT NOT NULL COMMENT 'id of the event table',
  user_id INT NOT NULL COMMENT 'user id of the person who is available',
  start_time DATETIME NOT NULL COMMENT 'start time of the availability',
  end_time DATETIME NOT NULL COMMENT 'end time of the availability',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (event_id) REFERENCES event_detail(id) ON DELETE CASCADE
  );