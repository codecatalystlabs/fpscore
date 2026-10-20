-- Individuals present, supervisors, and the action plan belong to a domain.
ALTER TABLE rhspars_attendees
    ADD COLUMN IF NOT EXISTS domain_id INTEGER REFERENCES rhspars_domains(id) ON DELETE SET NULL;

ALTER TABLE rhspars_action_plan_items
    ADD COLUMN IF NOT EXISTS domain_id INTEGER REFERENCES rhspars_domains(id) ON DELETE SET NULL;
