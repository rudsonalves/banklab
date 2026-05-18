DELETE FROM customer_documents cd
USING customers c
WHERE cd.customer_id = c.id
    AND cd.type = 'cpf'
    AND cd.value = c.cpf
    AND cd.country = 'BR';