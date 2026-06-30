--
-- PostgreSQL database dump
--

\restrict 1bNUaojGmUKu0063fLZef2XZUfozSjiwANyqtqyC6MleLtLNxaDrSTTJnvIxUrc

-- Dumped from database version 16.13
-- Dumped by pg_dump version 16.13

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.users VALUES ('b1486155-eb83-46c5-9f77-2d1e7673d36d', 'gopher', '$argon2id$v=19$m=65536,t=1,p=4$gmqRhs4/Neoo708g05eW8A$fpWmdn+np0pi31xsvBqwEmcXaVnGGxwnD3NPnAHrm9k', '2026-03-01 00:00:00', '2026-03-01 00:00:00');


--
-- Data for Name: expense_categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.expense_categories VALUES ('dd654864-25bc-42db-87a2-c2158a49519b', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'House', 'needs', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.expense_categories VALUES ('b7cd2e91-f4fa-4359-8413-87f44e83063a', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Utilities', 'needs', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.expense_categories VALUES ('6833f201-1557-4615-b28b-0a928cc508d5', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Transportation', 'needs', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.expense_categories VALUES ('4024bbac-7642-431a-8924-2accdffb4fcb', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Groceries', 'needs', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.expense_categories VALUES ('71e3b7dc-1ce3-4c93-81ac-44cc8b1dd26f', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Entertainment', 'wants', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.expense_categories VALUES ('924fdfe0-d705-4fc0-9028-620f2f257572', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Clothing', 'wants', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.expense_categories VALUES ('344f4403-c3da-42b4-99e8-ff639add0d01', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Shopping', 'wants', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.expense_categories VALUES ('5771af42-b305-4742-84e8-b8111fcd2d46', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Eating Out', 'wants', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.expense_categories VALUES ('9d75d005-4ee1-4ca1-941b-a4b56e480500', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Health', 'needs', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.expense_categories VALUES ('c57cd9f7-46a7-4b81-9e45-872200995d00', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Household Items', 'needs', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.expense_categories VALUES ('1d13701a-b78f-4201-bb22-07a35a63e7c6', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Personal', 'wants', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.expense_categories VALUES ('16777ce0-6bc6-4a09-8880-5c82892a9700', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Savings', 'savings', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.expense_categories VALUES ('6a8cf86c-8648-40fc-b972-a001191ea664', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Investing', 'savings', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.expense_categories VALUES ('61d859c8-502d-499c-998b-7969215fb367', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Emergency Fund', 'savings', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.expense_categories VALUES ('bc8b13d1-3643-40b4-8c13-4dc46c26b5e3', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Gifts/Donations', 'wants', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.expense_categories VALUES ('7028224b-1a8a-47b2-b810-e00f8e70124e', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Miscellaneous', 'needs', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);


--
-- Data for Name: accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.accounts VALUES ('61a9f1db-6eda-424c-afb3-0a088df81d29', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Main', 'bank', false, 'EUR', NULL, NULL, NULL, NULL, '2026-05-25 18:48:55.826741', '2026-05-25 18:48:55.826741', NULL);
INSERT INTO public.accounts VALUES ('a1f9fd71-861b-4e9d-82cc-acc56f6d9024', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Cash', 'bank', false, 'EUR', NULL, NULL, NULL, NULL, '2026-05-25 18:49:19.310473', '2026-05-25 18:49:19.310473', NULL);
INSERT INTO public.accounts VALUES ('a7b6fe5b-5820-4bdb-9d87-b76f2b69c46e', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Savings', 'bank', true, 'EUR', NULL, NULL, NULL, NULL, '2026-05-25 18:49:31.509823', '2026-05-25 18:49:31.509823', NULL);
INSERT INTO public.accounts VALUES ('d7158788-5015-43e1-a617-d9e5f2ad1e72', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Holidays', 'goal', NULL, 'EUR', 50000, '2026-05-01 00:00:00', '2026-12-01 00:00:00', '71e3b7dc-1ce3-4c93-81ac-44cc8b1dd26f', '2026-05-25 18:51:10.893844', '2026-05-25 18:51:10.893844', NULL);


--
-- Data for Name: income_categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.income_categories VALUES ('5d73ddb6-f839-4b83-8457-ac238ca66bec', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Job', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.income_categories VALUES ('694d93cc-fec8-4073-a361-850f6176866c', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Side Hustle', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.income_categories VALUES ('acb92080-2ba7-43ed-9d72-5466c77f725f', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Gifts', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);
INSERT INTO public.income_categories VALUES ('5b973f89-5440-4edc-9a1d-57d73f03d982', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'Others', '2026-03-11 00:00:00', '2026-03-11 00:00:00', NULL);


--
-- Data for Name: monthly_budgets; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.monthly_budgets VALUES ('30d03799-1930-46aa-8de0-e0f2c0991452', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', '4024bbac-7642-431a-8924-2accdffb4fcb', '2026-05-25', 25000, '2026-05-25 18:57:17.166853', '2026-05-25 18:57:17.166853');
INSERT INTO public.monthly_budgets VALUES ('c1ea41e4-63fe-4c4c-8c94-e183df2934d4', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', '6833f201-1557-4615-b28b-0a928cc508d5', '2026-05-25', 3000, '2026-05-25 18:57:29.396574', '2026-05-25 18:58:00.704453');
INSERT INTO public.monthly_budgets VALUES ('f9510e38-cbdc-496b-8b23-e25018c960d6', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'b7cd2e91-f4fa-4359-8413-87f44e83063a', '2026-05-25', 15000, '2026-05-25 18:59:08.519723', '2026-05-25 18:59:45.556278');
INSERT INTO public.monthly_budgets VALUES ('fb27ebbe-3053-4a2f-a432-487da0747ac0', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'dd654864-25bc-42db-87a2-c2158a49519b', '2026-05-25', 60000, '2026-05-25 19:00:23.761825', '2026-05-25 19:00:38.155267');
INSERT INTO public.monthly_budgets VALUES ('42d78411-b063-4ea3-a789-faba1d978f1b', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', '16777ce0-6bc6-4a09-8880-5c82892a9700', '2026-05-25', 10000, '2026-05-25 19:01:12.503269', '2026-05-25 19:01:35.104947');
INSERT INTO public.monthly_budgets VALUES ('6f8a1e4e-8b2c-4b9d-a005-109fe25a5c4e', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', '71e3b7dc-1ce3-4c93-81ac-44cc8b1dd26f', '2026-05-25', 7000, '2026-05-25 19:02:08.359642', '2026-05-25 19:02:18.053866');


--
-- Data for Name: ledger; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.ledger VALUES ('91a0b645-e04d-4ee0-9bea-2bac679e3856', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'expense', '61a9f1db-6eda-424c-afb3-0a088df81d29', NULL, NULL, '4024bbac-7642-431a-8924-2accdffb4fcb', 7000, 'Week groceries', '2026-05-03 00:00:00', '2026-05-25 19:06:52.509598');
INSERT INTO public.ledger VALUES ('7acc7ab5-a778-433e-9359-9ae4d22361af', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'income', NULL, '61a9f1db-6eda-424c-afb3-0a088df81d29', '5d73ddb6-f839-4b83-8457-ac238ca66bec', NULL, 100000, 'Salary', '2026-05-01 00:00:00', '2026-05-25 18:52:18.827814');
INSERT INTO public.ledger VALUES ('4670003b-a9d1-4944-9c4b-32d0daa78dee', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'expense', '61a9f1db-6eda-424c-afb3-0a088df81d29', NULL, NULL, 'dd654864-25bc-42db-87a2-c2158a49519b', 60000, 'Rent', '2026-05-04 00:00:00', '2026-05-25 19:16:59.442573');
INSERT INTO public.ledger VALUES ('49e61fc6-f475-42c1-9e09-dd10fc25b58d', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'expense', '61a9f1db-6eda-424c-afb3-0a088df81d29', NULL, NULL, '6833f201-1557-4615-b28b-0a928cc508d5', 3000, 'Train pass', '2026-05-01 00:00:00', '2026-05-25 19:19:44.525785');
INSERT INTO public.ledger VALUES ('8f8ed96f-6db6-4748-85f7-16706a5163f4', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'income', NULL, 'a1f9fd71-861b-4e9d-82cc-acc56f6d9024', '694d93cc-fec8-4073-a361-850f6176866c', NULL, 20000, 'Decluttering sale', '2026-05-02 00:00:00', '2026-05-25 18:55:13.113264');
INSERT INTO public.ledger VALUES ('7c7029a2-281c-4c0f-a848-a2798c49c940', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'transfer', 'a1f9fd71-861b-4e9d-82cc-acc56f6d9024', 'd7158788-5015-43e1-a617-d9e5f2ad1e72', NULL, '71e3b7dc-1ce3-4c93-81ac-44cc8b1dd26f', 7000, 'Holidays contribution', '2026-05-03 00:00:00', '2026-05-25 19:14:09.076749');
INSERT INTO public.ledger VALUES ('622d13cf-9d8a-42bd-8219-20f007525721', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'expense', '61a9f1db-6eda-424c-afb3-0a088df81d29', NULL, NULL, '4024bbac-7642-431a-8924-2accdffb4fcb', 5000, 'Week groceries', '2026-05-11 00:00:00', '2026-05-25 19:27:41.076994');
INSERT INTO public.ledger VALUES ('b8afd126-03d1-42b3-acc3-32abc04922d8', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'expense', '61a9f1db-6eda-424c-afb3-0a088df81d29', NULL, NULL, '4024bbac-7642-431a-8924-2accdffb4fcb', 5000, 'Week groceries', '2026-05-18 00:00:00', '2026-05-25 19:28:01.792538');
INSERT INTO public.ledger VALUES ('f79f7d9f-7824-4374-9d5e-64bbd4b83445', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'expense', '61a9f1db-6eda-424c-afb3-0a088df81d29', NULL, NULL, '4024bbac-7642-431a-8924-2accdffb4fcb', 5000, 'Week groceries', '2026-05-25 00:00:00', '2026-05-25 19:28:22.585553');
INSERT INTO public.ledger VALUES ('6e9f21b7-b34e-440a-a49e-8c7b1ea68fb8', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'expense', '61a9f1db-6eda-424c-afb3-0a088df81d29', NULL, NULL, 'b7cd2e91-f4fa-4359-8413-87f44e83063a', 15500, 'Utilities', '2026-05-05 00:00:00', '2026-05-25 19:30:20.4581');
INSERT INTO public.ledger VALUES ('b48b9fa6-2e66-42e4-9b14-337e1706b932', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'transfer', '61a9f1db-6eda-424c-afb3-0a088df81d29', 'a7b6fe5b-5820-4bdb-9d87-b76f2b69c46e', NULL, '16777ce0-6bc6-4a09-8880-5c82892a9700', 12500, 'Monthly Savings', '2026-05-31 00:00:00', '2026-05-25 19:32:57.292487');
INSERT INTO public.ledger VALUES ('e8456fdd-7402-4adc-9526-d1ba163da2d4', 'b1486155-eb83-46c5-9f77-2d1e7673d36d', 'transfer', 'a1f9fd71-861b-4e9d-82cc-acc56f6d9024', '61a9f1db-6eda-424c-afb3-0a088df81d29', NULL, NULL, 13000, 'Deposit', '2026-05-20 00:00:00', '2026-05-25 19:34:36.659196');


--
-- PostgreSQL database dump complete
--

\unrestrict 1bNUaojGmUKu0063fLZef2XZUfozSjiwANyqtqyC6MleLtLNxaDrSTTJnvIxUrc
