CREATE TYPE company_type AS ENUM ('Corporations','NonProfit','Cooperative','Sole Proprietorship');

CREATE TABLE "company"
(
    "id"                    uuid PRIMARY KEY,
    "name"                  varchar(15)   NOT NULL UNIQUE,
    "description"           varchar(3000),
    "amount_of_employees"   int NOT NULL,
    "registered"            bool NOT NULL,
    "type"          company_type NOT NULL
);
