-- 1. Categorías de contenido (opcional pero recomendado para filtrar noticias)
CREATE TABLE IF NOT EXISTS article_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Tabla Principal de Artículos / Noticias
CREATE TABLE IF NOT EXISTS articles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    publisher_id UUID NOT NULL REFERENCES publishers(id) ON DELETE CASCADE,
    category_id UUID NULL REFERENCES article_categories(id) ON DELETE SET NULL,
    
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE, -- Para URLs amigables (ej: /news/titanes-ganan-la-final)
    summary TEXT NOT NULL,              -- El extracto corto ("Leer más")
    content_body TEXT NOT NULL,         -- El texto completo / HTML / Markdown
    cover_image_url TEXT NULL,          -- Imagen principal de portada
    cloudinary_public_id VARCHAR(255) NULL,
    
    status VARCHAR(20) NOT NULL DEFAULT 'published', -- 'draft', 'published', 'archived'
    is_featured BOOLEAN DEFAULT false,              -- Para destacar en la home o carousel
    views_count INT DEFAULT 0,                      -- Contador simple de lecturas acumuladas
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    published_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID NOT NULL,
    name_created_by VARCHAR(100) NOT NULL
);

-- 3. Índices para Consultas Frecuentes
CREATE INDEX idx_articles_publisher ON articles (publisher_id);
CREATE INDEX idx_articles_status_published ON articles (status, published_at DESC);
CREATE INDEX idx_articles_slug ON articles (slug);