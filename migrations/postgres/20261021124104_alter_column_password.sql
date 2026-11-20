-- +goose Up
-- Сначала создаём новый частичный уникальный индекс.
CREATE UNIQUE INDEX users_email_active_uidx
    ON public.users (email)
    WHERE status <> 'inactive';

-- Удаляем глобальную уникальность email.
-- Индекс, принадлежащий этому constraint, удалится вместе с ним.
ALTER TABLE public.users
    DROP CONSTRAINT users_email_key;

-- Только если это отдельный, дополнительно созданный индекс.
DROP INDEX IF EXISTS public.users_email_idx;

-- +goose Down
ALTER TABLE public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);

CREATE INDEX users_email_idx
    ON public.users (email);

DROP INDEX public.users_email_active_uidx;
