# Sports Session Manager

A backend API for managing casual football and badminton sessions for a group of 50 players. Built to replace manual WhatsApp coordination with a structured system that enforces group rules automatically.

## The Problem

Every week, one organiser manually manages session scheduling, player confirmations, dropouts, penalties, and team generation over WhatsApp. No audit trail, frequent disputes about slot order, and the organiser spends 20-30 minutes per session on coordination overhead.

## What This System Does

- Players join sessions — goalkeeper priority enforced automatically
- Waitlist managed as a strict FIFO queue
- Match-day dropouts trigger financial penalties and strikes automatically
- 3 strikes boots a player from the group
- Outside players called only when group doesn't fill (70% threshold)
- Poll-based date voting when organiser is undecided
- Team generation with position balance and skill rating distribution

## Tech Stack

- **Language:** Go (Gin framework)
- **Database:** PostgreSQL
- **Architecture:** Handler → Service → Repository (3-layer)

## Project Structure

```
sports-session-manager/
├── config/         # Environment config
├── db/             # Database connection
├── handlers/       # HTTP handlers (player, session, team, penalty)
├── middleware/     # Request middleware
├── migrations/     # SQL migration files (001-006)
├── models/         # Data models and request/response types
├── repository/     # Database queries
├── services/       # Business logic
└── main.go
```

## Database Schema

6 tables:

| Table | Purpose |
|-------|---------|
| `players` | Group members and outside players with strike tracking |
| `sessions` | Match sessions with status lifecycle |
| `date_polls` | Poll options when organiser votes on date/format |
| `poll_votes` | Player votes on poll options |
| `session_slots` | Confirmed, waitlisted, and dropped slot records |
| `strike_events` | Audit trail for late dropouts and late arrivals |
| `penalties` | Outstanding ground fee payments |

## Business Rules Implemented

**Slot Priority (in order):**
1. Goalkeepers — first 2 GK slots always reserved
2. Group members — fill remaining slots
3. Outside players — called only when group can't fill (organiser notified at 70% threshold)
4. Outside players once confirmed cannot be displaced by group members

**Strike System:**
- Late dropout (match day) → strike + penalty
- Late arrival (organiser marks) → strike
- 3 strikes → player auto-booted, cannot join future sessions

**Dropout Logic (single transaction):**
- Slot marked as dropped
- If match day + was confirmed → penalty calculated (ground fee / confirmed count)
- Strike recorded
- First waitlisted player automatically promoted

**Team Balancing:**
- Players separated by position (GK, defender, midfielder, forward)
- Each group sorted by skill rating descending
- Snake draft within each position group
- Swap correction if skill difference > 2 after draft

## API Endpoints

### Players
```
POST   /players                    Create player
GET    /players                    List all players
GET    /players/:id                Get player by ID
PUT    /players/:id                Update player
POST   /players/:id/boot           Boot player (organiser)
POST   /players/:id/strike         Add strike manually
GET    /players/:id/strikes        Get strike history
GET    /players/:id/penalties      Get player's penalty history
```

### Sessions
```
POST   /sessions                   Create session (direct or poll)
GET    /sessions                   List all sessions
GET    /sessions/:id               Get session with slot details
POST   /sessions/:id/join          Join session
POST   /sessions/:id/drop          Drop from session
POST   /sessions/:id/confirm       Confirm session (organiser)
POST   /sessions/:id/complete      Complete session (organiser)
GET    /sessions/:id/check-outside Check if outside players should be called
GET    /sessions/:id/penalties     Get all penalties for a session
```

### Polls
```
POST   /sessions/:id/polls                    Create poll options
GET    /sessions/:id/polls                    Get polls with vote counts
POST   /sessions/:id/polls/:pollId/vote       Vote for a poll option
POST   /sessions/:id/polls/:pollId/lock       Lock winning option
```

### Teams
```
POST   /sessions/:id/generate-teams    Generate balanced teams
```

### Penalties
```
POST   /penalties/:id/mark-paid    Mark penalty as paid (organiser)
```

## Session Lifecycle

```
Direct session:
  created (open) → players join → confirmed → completed

Poll session:
  created (polling) → players vote → poll locked (open) → players join → confirmed → completed
```

## Setup

### Prerequisites
- Go 1.21+
- PostgreSQL

### Local Development

```bash
# Clone the repo
git clone https://github.com/mathewjacob123/sports-session-manager
cd sports-session-manager

# Set environment variables
cp .env.example .env
# Edit .env with your database credentials

# Run migrations
psql -d sports_session_manager -f migrations/001_create_players.sql
psql -d sports_session_manager -f migrations/002_create_sessions.sql
psql -d sports_session_manager -f migrations/003_create_date_polls.sql
psql -d sports_session_manager -f migrations/004_create_session_slots.sql
psql -d sports_session_manager -f migrations/005_create_strike_events.sql
psql -d sports_session_manager -f migrations/006_create_penalties.sql

# Run the server
go run main.go
```

### Environment Variables

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=sports_session_manager
PORT=8080
```

## Key Technical Decisions

**Why Go + Gin:** Gin is the most widely adopted Go web framework with excellent performance and clean middleware support. Chosen for its industry adoption and suitability for REST APIs.

**Why PostgreSQL:** Data is inherently relational — players, sessions, slots, and penalties all reference each other with foreign keys. PostgreSQL's RETURNING clause, CHECK constraints, and transaction support were used heavily throughout.

**Why layered architecture (Handler → Service → Repository):** Business logic lives in the service layer, database queries in the repository layer, and HTTP handling in the handler layer. Each layer has one responsibility. Changing the team balancing algorithm, for example, requires touching only `services/team_service.go`.

**Why transactions for dropout logic:** A match-day dropout touches 4 tables atomically — slot update, penalty creation, strike event, waitlist promotion. A transaction guarantees all four steps happen or none of them do. A failed waitlist promotion with a successful slot drop would leave the queue broken.

**Snake draft for team balancing:** Position balance is the primary constraint, skill balance is secondary. Once you enforce position grouping, the remaining optimisation space per group is small enough that snake draft and greedy converge. A swap correction step catches any skill gap over 2 points.

## Author

Mathew Jacob  
mathewjacob2503@gmail.com  
[LinkedIn](https://linkedin.com/in/mathew-jacob25)
