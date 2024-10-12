INSERT INTO users (email, username, password, role)
VALUES
    ('admin@example.com', 'admin', 'admin_password_hash', 'admin'),
    ('user1@example.com', 'user1', 'user1_password_hash', 'user'),
    ('user2@example.com', 'user2', 'user2_password_hash', 'user');

INSERT INTO documents (name, mime_type, file, public, location, token, owner_email)
VALUES
    ('photo.jpg', 'image/jpeg', true, false, '/uploads/photo.jpg', 'token123', 'admin@example.com'),
    ('document.pdf', 'application/pdf', true, true, '/uploads/document.pdf', 'token456', 'user1@example.com'),
    ('music.mp3', 'audio/mpeg', true, false, '/uploads/music.mp3', 'token789', 'user2@example.com');


INSERT INTO document_grants (document_id, user_email)
VALUES
    ((SELECT id FROM documents WHERE name = 'photo.jpg'), 'user1@example.com'),
    ((SELECT id FROM documents WHERE name = 'document.pdf'), 'user2@example.com');
