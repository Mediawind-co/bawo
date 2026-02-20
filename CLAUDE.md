# Bawo Project - Claude Code Memory

## Project Overview
Bawo is a language learning platform built with:
- **Backend**: Go with Encore.dev framework (v1.52.1)
- **Frontend**: Next.js (in `/web` directory)
- **Database**: PostgreSQL (via Encore)

## Project Structure
```
bawo/
├── admin/           # Admin panel services
│   ├── adminauth/   # Admin authentication & superadmin seeding
│   └── dashboard/   # Admin dashboard API
├── identity/        # User identity management
│   ├── auth/        # User authentication (Google, Apple OAuth)
│   └── user/        # User profiles & management
├── learning/        # Learning content services
│   ├── content/     # Units, lessons, questions, audio
│   ├── language/    # Supported languages
│   └── lesson/      # Lesson sessions & evaluation
├── progress/        # User progress tracking
│   ├── enrollment/  # Course enrollments
│   ├── streak/      # Learning streaks
│   └── tracker/     # Progress & mistakes tracking
└── web/             # Next.js frontend
```

## Superadmin Credentials
**File**: `admin/adminauth/seed.go`

| Field    | Value              |
|----------|-------------------|
| Username | `superadmin`      |
| Email    | `admin@bawo.app`  |
| Password | `BawoAdmin2024!`  |
| Name     | `Super Admin`     |

The superadmin is automatically seeded on service initialization via `seedSuperadmin()` function.

## API Configuration
- **Backend URL**: `http://localhost:4000`
- **Google Client ID**: `374880305181-3clq8r0vumm9128ua7vsutdg7p1cb8p9.apps.googleusercontent.com`

## Key Services

### Admin Auth (`admin/adminauth/`)
- Login endpoint: `/admin/auth/login`
- Token validation endpoint
- Session management (in-memory tokens with `admin_` prefix)

### User Auth (`identity/auth/`)
- OAuth providers: Google, Apple
- Development login support
- JWT token validation

### Content Management (`learning/content/`)
- Units, Lessons, Questions CRUD
- Audio file handling for language learning

## Database Tables
- `admins` - Admin users with superadmin flag
- `users` - Regular users
- `languages` - Supported languages
- `units`, `lessons`, `questions` - Learning content
- `enrollments` - User course enrollments
- `progress`, `mistakes` - User learning progress
- `streaks` - User learning streaks
- `sessions` - Lesson sessions

## Development Commands
```bash
# Run the backend
encore run

# Run the frontend
cd web && npm run dev
```

## Git Workflow
- **CI/CD**: Vercel (frontend) + Encore (backend) connected to GitHub
- **Branch strategy**: Use feature branches for new work, merge to main when complete
- Create branches like: `feature/admin-management`, `fix/auth-bug`, etc.
- **IMPORTANT**: Before starting work on any feature or fix, ALWAYS:
  1. Switch to main and pull latest changes
  2. Create a new branch from main
  ```bash
  git checkout main
  git pull origin main
  git checkout -b feature/your-feature-name
  # or
  git checkout -b fix/your-fix-name
  ```

## Work History

### Session: 2026-02-18
- Reviewed superadmin credentials and admin authentication system
- Created this CLAUDE.md file for project memory
- **Fixed admin auth context error**: Admin pages were using `useAuth` (user auth) instead of `useAdminAuth` (admin auth), causing "useAuth must be used within an AuthProvider" error on login
  - Fixed files: `web/src/app/admin/page.tsx`, `users/page.tsx`, `languages/page.tsx`, `content/page.tsx`
  - Changed imports from `@/lib/auth-context` to `@/lib/admin-auth-context`
- **Fixed admin token validation**: Added `admin_` token support to the main AuthHandler in `identity/auth/auth.go`
  - Admin tokens are now validated via `adminauth.GetAdminByToken()`
  - Returns auth data with `role: "admin"` and `provider: "admin"`
- **Fixed AdminOnly middleware role extraction**: Updated `getRoleFromAuthData()` in `identity/user/user.go`
  - Now uses reflection to properly extract the Role field from AuthData struct
  - Handles both pointer and value types
- **Added Admin User Management feature**:
  - Backend endpoints in `admin/adminauth/auth.go`:
    - GET `/admin/admins` - List all admins
    - GET `/admin/admins/:id` - Get specific admin
    - POST `/admin/admins` - Create new admin
    - PUT `/admin/admins/:id` - Update admin
    - DELETE `/admin/admins/:id` - Delete admin (prevents deleting last superadmin)
    - PUT `/admin/admins/:id/password` - Change admin password
  - Frontend page at `web/src/app/admin/admins/page.tsx`
  - Added ShieldIcon to icons component
  - Added "Admins" navigation link in admin layout

### Session: 2026-02-20
- **Fixed 404 on learn unit/lesson pages**: Created missing dynamic route pages
  - Added `web/src/app/(dashboard)/learn/[languageId]/unit/[unitId]/page.tsx` - Unit page showing lessons
  - Added `web/src/app/(dashboard)/learn/[languageId]/unit/[unitId]/lesson/[lessonId]/page.tsx` - Lesson taking page
  - Added `ArrowLeftIcon` and `StarIcon` to icons component
- **Fixed mistake tracking on wrong answers**: Mistakes were not being saved when users answered incorrectly
  - Added `GetLanguageIDByQuestion()` helper in `learning/content/repository.go` to get language ID from question
  - Added `MistakeRecordedEvent` and `MistakeRecorded` pub/sub topic in `learning/lesson/events.go`
  - Updated `SubmitAnswer()` in `learning/lesson/lesson.go` to publish mistake event on wrong answers
  - Added `HandleMistakeRecorded()` subscriber in `progress/tracker/tracker.go` to save mistakes to database

## Important Notes
- **Two separate auth systems**:
  - User auth: `useAuth` from `@/lib/auth-context` (uses `bawo_token`)
  - Admin auth: `useAdminAuth` from `@/lib/admin-auth-context` (uses `bawo_admin_token`)
- Admin pages must use `useAdminAuth`, not `useAuth`
