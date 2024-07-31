-- +goose Up
INSERT INTO users (username) VALUES
    ('user001'),
    ('user002'),
    ('user003'),
    ('user004'),
    ('user005'),
    ('user006'),
    ('user007'),
    ('user008'),
    ('user009'),
    ('user010'),
    ('user011');

-- Inserting scenarios
INSERT INTO scenarios (name, description) VALUES
    ('Scenario 1', 'Description for Scenario 1'),
    ('Scenario 2', 'Description for Scenario 2'),
    ('Scenario 3', 'Description for Scenario 3'),
    ('Scenario 4', 'Description for Scenario 4'),
    ('Scenario 5', 'Description for Scenario 5'),
    ('Scenario 6', 'Description for Scenario 6'),
    ('Scenario 7', 'Description for Scenario 7'),
    ('Scenario 8', 'Description for Scenario 8'),
    ('Scenario 9', 'Description for Scenario 9'),
    ('Scenario 10', 'Description for Scenario 10');



-- Inserting blocks
INSERT INTO blocks (key, title, background_image) VALUES
    ('block001', 'Block 1', 'image1.jpg'),
    ('block002', 'Block 2', 'image2.jpg'),
    ('block003', 'Block 3', NULL),
    ('block004', 'Block 4', 'image4.jpg'),
    ('block005', 'Block 5', 'image5.jpg'),
    ('block006', 'Block 6', 'image6.jpg'),
    ('block007', 'Block 7', 'image7.jpg'),
    ('block008', 'Block 8', NULL),
    ('block009', 'Block 9', 'image9.jpg'),
    ('block010', 'Block 10', 'image10.jpg'),
    ('block011', 'Block 11', NULL),
    ('block012', 'Block 12', 'image12.jpg'),
    ('block013', 'Block 13', 'image13.jpg'),
    ('block014', 'Block 14', NULL),
    ('block015', 'Block 15', 'image15.jpg'),
    ('block016', 'Block 16', 'image16.jpg');

-- Inserting scenario_user mappings
INSERT INTO scenario_user (scenario_id, user_id) VALUES
    (1, 1),   -- Scenario 1 with User 1
    (1, 2),   -- Scenario 1 with User 2
    (2, 1),   -- Scenario 2 with User 1
    (3, 3),   -- Scenario 3 with User 3
    (4, 4),   -- Scenario 4 with User 4
    (5, 5),   -- Scenario 5 with User 5
    (6, 6),   -- Scenario 6 with User 6
    (6, 7),   -- Scenario 6 with User 7
    (7, 8),   -- Scenario 7 with User 8
    (8, 9),   -- Scenario 8 with User 9
    (9, 10),  -- Scenario 9 with User 10
    (10, 1),  -- Scenario 10 with User 1
    (10, 2);  -- Scenario 10 with User 2

-- Inserting scenario_mapping
INSERT INTO scenario_mapping (scenario_id, key) VALUES
    (1, 'block001'),
    (1, 'block002'),
    (2, 'block003'),
    (3, 'block004'),
    (4, 'block005'),
    (6, 'block006'),
    (6, 'block007'),
    (7, 'block008'),
    (8, 'block009'),
    (9, 'block010'),
    (10, 'block001'),
    (10, 'block002'),
    (10, 'block003'),
    (10, 'block004'),
    (10, 'block005');

-- Inserting content_type_dic
INSERT INTO content_type_dic (name, description) VALUES
    ('Text', 'Plain text content'),
    ('Image', 'Image content'),
    ('Video', 'Video content'),
    ('Audio', 'Audio content'),
    ('Document', 'Document content'),
    ('Link', 'Link to external content'),
    ('Code', 'Code snippet'),
    ('Quiz', 'Interactive quiz');

-- Inserting contents
INSERT INTO contents (name, content_type_id, content) VALUES
    ('Content 1', 1, '{"text": "Sample text content for Block 1"}'),
    ('Content 2', 2, '{"url": "image-url.jpg"}'),
    ('Content 3', 3, '{"url": "video-url.mp4"}'),
    ('Content 4', 4, '{"url": "audio-url.mp3"}'),
    ('Content 5', 5, '{"url": "document-url.pdf"}'),
    ('Content 6', 6, '{"url": "external-link-url"}'),
    ('Content 7', 7, '{"code": "function() { return true; }"}'),
    ('Content 8', 8, '{"question": "What is 1+1?", "options": ["1", "2", "3"], "correct_option": 2}'),
    ('Content 9', 1, '{"text": "Sample text content for Block 3"}'),
    ('Content 10', 2, '{"url": "image-url-2.jpg"}'),
    ('Content 11', 3, '{"url": "video-url-2.mp4"}'),
    ('Content 12', 4, '{"url": "audio-url-2.mp3"}'),
    ('Content 13', 5, '{"url": "document-url-2.pdf"}'),
    ('Content 14', 6, '{"url": "external-link-url-2"}'),
    ('Content 15', 7, '{"code": "function() { return false; }"}'),
    ('Content 16', 8, '{"question": "What is 2+2?", "options": ["2", "3", "4"], "correct_option": 3}');

-- Inserting content_mapping
INSERT INTO content_mapping (content_block_id, content_id, rating) VALUES
    (1, 1, 5),    -- Block 1 with Content 1, rating 5
    (2, 2, 4),    -- Block 2 with Content 2, rating 4
    (3, 3, NULL), -- Block 3 with Content 3, no rating
    (4, 4, 3),    -- Block 4 with Content 4, rating 3
    (5, 5, 2),    -- Block 5 with Content 5, rating 2
    (6, 6, NULL), -- Block 6 with Content 6
    (7, 7, NULL), -- Block 7 with Content 7
    (8, 8, NULL), -- Block 8 with Content 8
    (9, 9, NULL), -- Block 9 with Content 9
    (10, 10, NULL), -- Block 10 with Content 10
    (11, 11, NULL), -- Block 11 with Content 11
    (12, 12, NULL), -- Block 12 with Content 12
    (13, 13, NULL), -- Block 13 with Content 13
    (14, 14, NULL), -- Block 14 with Content 14
    (15, 15, NULL), -- Block 15 with Content 15
    (16, 16, NULL); -- Block 16 with Content 16

-- Inserting template_contents
INSERT INTO template_contents (name, template_content, content_type_id) VALUES
    ('Template 1', '{"content": "Template content 1"}', 1),
    ('Template 2', '{"content": "Template content 2"}', 2),
    ('Template 3', '{"content": "Template content 3"}', 3),
    ('Template 4', '{"content": "Template content 4"}', 4),
    ('Template 5', '{"content": "Template content 5"}', 5),
    ('Template 6', '{"content": "Template content 6"}', 6),
    ('Template 7', '{"content": "Template content 7"}', 7),
    ('Template 8', '{"content": "Template content 8"}', 8);

INSERT INTO content_mapping (content_block_id, content_id, rating) VALUES
    (1, 2, 4), 
    (1, 3, NULL);

-- +goose Down
