CREATE TABLE IF NOT EXISTS event_slot (
  id INT PRIMARY KEY AUTO_INCREMENT,
  event_id INT NOT NULL COMMENT 'id of the event table',
  start_time DATETIME NOT NULL COMMENT 'start time of the slot',
  end_time DATETIME NOT NULL COMMENT 'end time of the slot',
  FOREIGN KEY (event_id) REFERENCES event_detail(id) ON DELETE CASCADE ON UPDATE CASCADE
);