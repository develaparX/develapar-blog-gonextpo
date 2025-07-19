-- ========================================
-- IMPROVED DDL: CREATE TABLE
-- Versi yang diperbaiki dengan best practices
-- ========================================

-- Tabel users - Diperbaiki dengan kolom tambahan dan constraint yang lebih baik
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,           -- Tambah username untuk login alternatif
  email VARCHAR(255) UNIQUE NOT NULL,             -- Perbesar ukuran email
  email_verified_at TIMESTAMPTZ,                  -- Verifikasi email
  first_name VARCHAR(100) NOT NULL,               -- Pisah nama depan dan belakang
  last_name VARCHAR(100),                         -- Nama belakang opsional
  display_name VARCHAR(200),                      -- Nama untuk ditampilkan
  password_hash VARCHAR(255) NOT NULL,            -- Lebih eksplisit bahwa ini hash
  avatar_url VARCHAR(500),                        -- URL avatar user
  bio TEXT,                                       -- Bio singkat user
  role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'editor', 'user')), -- Enum constraint
  status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')), -- Status user
  last_login_at TIMESTAMPTZ,                      -- Track login terakhir
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tabel categories - Diperbaiki dengan slug dan metadata
CREATE TABLE categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) UNIQUE NOT NULL,              -- Perbesar ukuran nama
  slug VARCHAR(100) UNIQUE NOT NULL,              -- Tambah slug untuk URL
  description TEXT,                               -- Deskripsi kategori
  color_hex VARCHAR(7),                           -- Warna kategori (#FFFFFF)
  icon_name VARCHAR(50),                          -- Nama icon untuk UI
  sort_order INTEGER DEFAULT 0,                  -- Urutan tampilan
  is_active BOOLEAN NOT NULL DEFAULT true,       -- Status aktif/nonaktif
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tabel tags - Diperbaiki dengan slug dan metadata
CREATE TABLE tags (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) UNIQUE NOT NULL,              -- Perbesar ukuran nama
  slug VARCHAR(100) UNIQUE NOT NULL,              -- Tambah slug untuk URL
  description TEXT,                               -- Deskripsi tag
  color_hex VARCHAR(7),                           -- Warna tag
  usage_count INTEGER NOT NULL DEFAULT 0,        -- Hitung penggunaan tag
  is_active BOOLEAN NOT NULL DEFAULT true,       -- Status aktif/nonaktif
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tabel articles - Diperbaiki dengan metadata lengkap
CREATE TABLE articles (
  id SERIAL PRIMARY KEY,
  title VARCHAR(500) NOT NULL,                    -- Perbesar ukuran judul
  slug VARCHAR(500) UNIQUE NOT NULL,              -- Perbesar ukuran slug
  excerpt TEXT,                                   -- Ringkasan artikel
  content TEXT NOT NULL,
  featured_image_url VARCHAR(500),               -- URL gambar utama
  meta_title VARCHAR(500),                        -- SEO meta title
  meta_description TEXT,                          -- SEO meta description
  author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- Lebih eksplisit
  category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived', 'deleted')), -- Status artikel
  visibility VARCHAR(20) NOT NULL DEFAULT 'public' CHECK (visibility IN ('public', 'private', 'password_protected')), -- Visibilitas
  password_hash VARCHAR(255),                     -- Password untuk artikel protected
  is_featured BOOLEAN NOT NULL DEFAULT false,    -- Artikel unggulan
  allow_comments BOOLEAN NOT NULL DEFAULT true,  -- Izinkan komentar
  view_count INTEGER NOT NULL DEFAULT 0,         -- Lebih eksplisit
  like_count INTEGER NOT NULL DEFAULT 0,         -- Cache count untuk performa
  comment_count INTEGER NOT NULL DEFAULT 0,      -- Cache count untuk performa
  reading_time_minutes INTEGER,                  -- Estimasi waktu baca
  published_at TIMESTAMPTZ,                      -- Waktu publikasi
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  
  -- Constraints
  CONSTRAINT check_password_for_protected CHECK (
    (visibility != 'password_protected') OR (password_hash IS NOT NULL)
  )
);

-- Tabel article_tags - Diperbaiki dengan metadata
CREATE TABLE article_tags (
  article_id INTEGER NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
  tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),  -- Track kapan tag ditambahkan
  created_by INTEGER REFERENCES users(id),        -- Siapa yang menambahkan tag
  PRIMARY KEY (article_id, tag_id)
);

-- Tabel comments - Diperbaiki dengan fitur lengkap
CREATE TABLE comments (
  id SERIAL PRIMARY KEY,
  article_id INTEGER NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
  author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- Lebih eksplisit
  parent_comment_id INTEGER REFERENCES comments(id) ON DELETE CASCADE, -- Nested comments
  content TEXT NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'approved' CHECK (status IN ('pending', 'approved', 'rejected', 'spam')), -- Moderasi
  like_count INTEGER NOT NULL DEFAULT 0,          -- Cache count
  is_edited BOOLEAN NOT NULL DEFAULT false,      -- Track apakah sudah diedit
  edited_at TIMESTAMPTZ,                          -- Waktu edit terakhir
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tabel comment_likes - Terpisah untuk like komentar
CREATE TABLE comment_likes (
  id SERIAL PRIMARY KEY,
  comment_id INTEGER NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (comment_id, user_id)
);

-- Tabel article_likes - Rename dari 'likes' untuk lebih spesifik
CREATE TABLE article_likes (
  id SERIAL PRIMARY KEY,
  article_id INTEGER NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (article_id, user_id)
);

-- Tabel bookmarks - Diperbaiki dengan metadata
CREATE TABLE bookmarks (
  id SERIAL PRIMARY KEY,
  article_id INTEGER NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  folder_name VARCHAR(100),                       -- Organisasi bookmark dalam folder
  notes TEXT,                                     -- Catatan pribadi untuk bookmark
  is_favorite BOOLEAN NOT NULL DEFAULT false,    -- Bookmark favorit
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (article_id, user_id)
);

-- Tabel refresh_tokens - Diperbaiki dengan security
CREATE TABLE refresh_tokens (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash VARCHAR(255) NOT NULL UNIQUE,       -- Hash token untuk security
  device_info JSONB,                             -- Info device (browser, OS, dll)
  ip_address INET,                               -- IP address untuk security
  is_revoked BOOLEAN NOT NULL DEFAULT false,     -- Status revoke
  expires_at TIMESTAMPTZ NOT NULL,
  last_used_at TIMESTAMPTZ,                      -- Track penggunaan terakhir
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tabel baru: article_views - Track view detail
CREATE TABLE article_views (
  id SERIAL PRIMARY KEY,
  article_id INTEGER NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
  user_id INTEGER REFERENCES users(id) ON DELETE SET NULL, -- Null untuk anonymous
  ip_address INET,                               -- IP untuk anonymous tracking
  user_agent TEXT,                               -- Browser info
  referrer_url VARCHAR(500),                     -- Dari mana datang
  view_duration_seconds INTEGER,                 -- Berapa lama membaca
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tabel baru: user_follows - Follow antar user
CREATE TABLE user_follows (
  id SERIAL PRIMARY KEY,
  follower_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  following_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (follower_id, following_id),
  CONSTRAINT no_self_follow CHECK (follower_id != following_id)
);

-- Tabel baru: notifications - Sistem notifikasi
CREATE TABLE notifications (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type VARCHAR(50) NOT NULL,                     -- 'comment', 'like', 'follow', dll
  title VARCHAR(200) NOT NULL,
  message TEXT NOT NULL,
  data JSONB,                                    -- Data tambahan (article_id, dll)
  is_read BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ========================================
-- INDEXES untuk Performance
-- ========================================

-- Users indexes
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- Categories indexes
CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_active ON categories(is_active);
CREATE INDEX idx_categories_sort_order ON categories(sort_order);

-- Tags indexes
CREATE INDEX idx_tags_slug ON tags(slug);
CREATE INDEX idx_tags_active ON tags(is_active);
CREATE INDEX idx_tags_usage_count ON tags(usage_count DESC);

-- Articles indexes
CREATE INDEX idx_articles_slug ON articles(slug);
CREATE INDEX idx_articles_author_id ON articles(author_id);
CREATE INDEX idx_articles_category_id ON articles(category_id);
CREATE INDEX idx_articles_status ON articles(status);
CREATE INDEX idx_articles_published_at ON articles(published_at DESC);
CREATE INDEX idx_articles_featured ON articles(is_featured);
CREATE INDEX idx_articles_view_count ON articles(view_count DESC);
CREATE INDEX idx_articles_created_at ON articles(created_at DESC);

-- Composite indexes untuk query kompleks
CREATE INDEX idx_articles_status_published ON articles(status, published_at DESC) WHERE status = 'published';
CREATE INDEX idx_articles_author_status ON articles(author_id, status);
CREATE INDEX idx_articles_category_status ON articles(category_id, status) WHERE status = 'published';

-- Comments indexes
CREATE INDEX idx_comments_article_id ON comments(article_id);
CREATE INDEX idx_comments_author_id ON comments(author_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_comment_id);
CREATE INDEX idx_comments_status ON comments(status);
CREATE INDEX idx_comments_created_at ON comments(created_at DESC);

-- Likes indexes
CREATE INDEX idx_article_likes_article_id ON article_likes(article_id);
CREATE INDEX idx_article_likes_user_id ON article_likes(user_id);
CREATE INDEX idx_comment_likes_comment_id ON comment_likes(comment_id);

-- Bookmarks indexes
CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX idx_bookmarks_folder ON bookmarks(folder_name);
CREATE INDEX idx_bookmarks_favorite ON bookmarks(is_favorite);

-- Article tags indexes
CREATE INDEX idx_article_tags_article_id ON article_tags(article_id);
CREATE INDEX idx_article_tags_tag_id ON article_tags(tag_id);

-- Views indexes
CREATE INDEX idx_article_views_article_id ON article_views(article_id);
CREATE INDEX idx_article_views_user_id ON article_views(user_id);
CREATE INDEX idx_article_views_created_at ON article_views(created_at DESC);

-- Notifications indexes
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_unread ON notifications(user_id, is_read) WHERE is_read = false;

-- Full-text search indexes
CREATE INDEX idx_articles_search ON articles USING gin(to_tsvector('english', title || ' ' || coalesce(excerpt, '') || ' ' || content));
CREATE INDEX idx_users_search ON users USING gin(to_tsvector('english', first_name || ' ' || coalesce(last_name, '') || ' ' || coalesce(display_name, '')));

-- ========================================
-- TRIGGERS untuk Auto-update
-- ========================================

-- Function untuk update timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers untuk auto-update updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tags_updated_at BEFORE UPDATE ON tags FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_articles_updated_at BEFORE UPDATE ON articles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_comments_updated_at BEFORE UPDATE ON comments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_bookmarks_updated_at BEFORE UPDATE ON bookmarks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function untuk update counter
CREATE OR REPLACE FUNCTION update_article_counters()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Update counter saat ada data baru
        IF TG_TABLE_NAME = 'article_likes' THEN
            UPDATE articles SET like_count = like_count + 1 WHERE id = NEW.article_id;
        ELSIF TG_TABLE_NAME = 'comments' THEN
            UPDATE articles SET comment_count = comment_count + 1 WHERE id = NEW.article_id;
        ELSIF TG_TABLE_NAME = 'article_views' THEN
            UPDATE articles SET view_count = view_count + 1 WHERE id = NEW.article_id;
        END IF;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        -- Update counter saat data dihapus
        IF TG_TABLE_NAME = 'article_likes' THEN
            UPDATE articles SET like_count = like_count - 1 WHERE id = OLD.article_id;
        ELSIF TG_TABLE_NAME = 'comments' THEN
            UPDATE articles SET comment_count = comment_count - 1 WHERE id = OLD.article_id;
        END IF;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ language 'plpgsql';

-- Triggers untuk auto-update counters
CREATE TRIGGER update_article_like_count AFTER INSERT OR DELETE ON article_likes FOR EACH ROW EXECUTE FUNCTION update_article_counters();
CREATE TRIGGER update_article_comment_count AFTER INSERT OR DELETE ON comments FOR EACH ROW EXECUTE FUNCTION update_article_counters();
CREATE TRIGGER update_article_view_count AFTER INSERT ON article_views FOR EACH ROW EXECUTE FUNCTION update_article_counters();

-- Function untuk update tag usage count
CREATE OR REPLACE FUNCTION update_tag_usage_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE tags SET usage_count = usage_count + 1 WHERE id = NEW.tag_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE tags SET usage_count = usage_count - 1 WHERE id = OLD.tag_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ language 'plpgsql';

-- Trigger untuk auto-update tag usage count
CREATE TRIGGER update_tag_usage_count AFTER INSERT OR DELETE ON article_tags FOR EACH ROW EXECUTE FUNCTION update_tag_usage_count();