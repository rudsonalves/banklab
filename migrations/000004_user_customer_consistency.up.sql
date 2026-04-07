ALTER TABLE users
    ADD CONSTRAINT chk_users_customer_role_consistency CHECK (
        (role = 'customer' AND customer_id IS NOT NULL)
        OR (role != 'customer')
    );
