-- +goose Up
-- +goose StatementBegin
INSERT INTO philosophers (name, date_born, date_died, birthplace, interests, portrait_uri, bio)
VALUES
  ('Socrates', '469 BCE', '399 BCE', 'Athens, Greece', '{"Ethics", "Epistemology", "Logic"}', NULL, 'Classical Greek philosopher who laid the groundwork for Western philosophy.'),
  ('Plato', '427 BCE', '347 BCE', 'Athens, Greece', '{"Metaphysics", "Ethics", "Politics"}', NULL, 'Student of Socrates and founder of the Academy. Known for the Theory of Forms.'),
  ('Aristotle', '384 BCE', '322 BCE', 'Stagira, Greece', '{"Logic", "Biology", "Ethics", "Politics"}', NULL, 'Student of Plato and tutor to Alexander the Great. Founder of the Peripatetic school.'),
  ('Immanuel Kant', '1724', '1804', 'Königsberg, Prussia', '{"Metaphysics", "Epistemology", "Ethics", "Aesthetics"}', NULL, 'Central figure in modern philosophy. Known for his work on transcendental idealism and moral duty.'),
  ('Friedrich Nietzsche', '1844', '1900', 'Röcken, Germany', '{"Existentialism", "Morality", "Religion", "Power"}', NULL, 'Influential critic of religion and morality. Introduced concepts like the Übermensch and will to power.');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM philosophers
WHERE name IN ('Socrates', 'Plato', 'Aristotle', 'Immanuel Kant', 'Friedrich Nietzsche');
-- +goose StatementEnd