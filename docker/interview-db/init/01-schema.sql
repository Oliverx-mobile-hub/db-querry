-- interview_db schema for MySQL 8.4
-- All names and personal details in this sample database are synthetic.

SET NAMES utf8mb4;
SET time_zone = '+08:00';

CREATE DATABASE IF NOT EXISTS interview_db
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;
USE interview_db;

CREATE TABLE offices (
  office_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  office_code VARCHAR(20) NOT NULL,
  office_name VARCHAR(100) NOT NULL,
  city VARCHAR(60) NOT NULL,
  address VARCHAR(255) NOT NULL,
  timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Shanghai',
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (office_id),
  UNIQUE KEY uk_offices_code (office_code)
) ENGINE=InnoDB COMMENT='办公地点';

CREATE TABLE departments (
  department_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  parent_department_id BIGINT UNSIGNED NULL,
  department_code VARCHAR(30) NOT NULL,
  department_name VARCHAR(100) NOT NULL,
  cost_center VARCHAR(30) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (department_id),
  UNIQUE KEY uk_departments_code (department_code),
  KEY idx_departments_parent (parent_department_id),
  CONSTRAINT fk_departments_parent FOREIGN KEY (parent_department_id)
    REFERENCES departments (department_id) ON DELETE SET NULL
) ENGINE=InnoDB COMMENT='组织部门';

CREATE TABLE employees (
  employee_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  employee_no VARCHAR(30) NOT NULL,
  full_name VARCHAR(80) NOT NULL,
  email VARCHAR(160) NOT NULL,
  department_id BIGINT UNSIGNED NOT NULL,
  office_id BIGINT UNSIGNED NOT NULL,
  job_title VARCHAR(100) NOT NULL,
  employee_role ENUM('recruiter','hiring_manager','interviewer','hrbp','approver') NOT NULL,
  hire_date DATE NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (employee_id),
  UNIQUE KEY uk_employees_no (employee_no),
  UNIQUE KEY uk_employees_email (email),
  KEY idx_employees_department (department_id),
  KEY idx_employees_office (office_id),
  CONSTRAINT fk_employees_department FOREIGN KEY (department_id)
    REFERENCES departments (department_id),
  CONSTRAINT fk_employees_office FOREIGN KEY (office_id)
    REFERENCES offices (office_id)
) ENGINE=InnoDB COMMENT='参与招聘的内部员工';

CREATE TABLE candidate_sources (
  source_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  source_code VARCHAR(30) NOT NULL,
  source_name VARCHAR(100) NOT NULL,
  source_category ENUM('career_site','referral','agency','campus','social','job_board','event') NOT NULL,
  cost_per_hire DECIMAL(12,2) NOT NULL DEFAULT 0,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  PRIMARY KEY (source_id),
  UNIQUE KEY uk_candidate_sources_code (source_code)
) ENGINE=InnoDB COMMENT='候选人来源渠道';

CREATE TABLE job_requisitions (
  requisition_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  requisition_no VARCHAR(30) NOT NULL,
  job_title VARCHAR(120) NOT NULL,
  department_id BIGINT UNSIGNED NOT NULL,
  office_id BIGINT UNSIGNED NOT NULL,
  hiring_manager_id BIGINT UNSIGNED NOT NULL,
  recruiter_id BIGINT UNSIGNED NOT NULL,
  employment_type ENUM('full_time','part_time','contract','intern') NOT NULL DEFAULT 'full_time',
  workplace_type ENUM('onsite','hybrid','remote') NOT NULL DEFAULT 'hybrid',
  job_level ENUM('intern','junior','mid','senior','lead','manager','director') NOT NULL,
  headcount INT UNSIGNED NOT NULL DEFAULT 1,
  filled_count INT UNSIGNED NOT NULL DEFAULT 0,
  min_salary DECIMAL(12,2) NOT NULL,
  max_salary DECIMAL(12,2) NOT NULL,
  currency CHAR(3) NOT NULL DEFAULT 'CNY',
  status ENUM('draft','open','on_hold','closed','cancelled') NOT NULL DEFAULT 'open',
  opened_at DATE NULL,
  target_close_date DATE NULL,
  closed_at DATE NULL,
  description TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (requisition_id),
  UNIQUE KEY uk_job_requisitions_no (requisition_no),
  KEY idx_requisitions_status_department (status, department_id),
  KEY idx_requisitions_recruiter (recruiter_id),
  CONSTRAINT chk_requisition_salary CHECK (min_salary > 0 AND max_salary >= min_salary),
  CONSTRAINT chk_requisition_headcount CHECK (filled_count <= headcount),
  CONSTRAINT fk_requisitions_department FOREIGN KEY (department_id)
    REFERENCES departments (department_id),
  CONSTRAINT fk_requisitions_office FOREIGN KEY (office_id)
    REFERENCES offices (office_id),
  CONSTRAINT fk_requisitions_manager FOREIGN KEY (hiring_manager_id)
    REFERENCES employees (employee_id),
  CONSTRAINT fk_requisitions_recruiter FOREIGN KEY (recruiter_id)
    REFERENCES employees (employee_id)
) ENGINE=InnoDB COMMENT='招聘需求';

CREATE TABLE skills (
  skill_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  skill_name VARCHAR(80) NOT NULL,
  skill_category ENUM('engineering','product','design','data','business','language','leadership') NOT NULL,
  PRIMARY KEY (skill_id),
  UNIQUE KEY uk_skills_name (skill_name)
) ENGINE=InnoDB COMMENT='技能字典';

CREATE TABLE requisition_skills (
  requisition_id BIGINT UNSIGNED NOT NULL,
  skill_id BIGINT UNSIGNED NOT NULL,
  importance ENUM('required','preferred') NOT NULL,
  minimum_years DECIMAL(3,1) NOT NULL DEFAULT 0,
  PRIMARY KEY (requisition_id, skill_id),
  CONSTRAINT fk_requisition_skills_requisition FOREIGN KEY (requisition_id)
    REFERENCES job_requisitions (requisition_id) ON DELETE CASCADE,
  CONSTRAINT fk_requisition_skills_skill FOREIGN KEY (skill_id)
    REFERENCES skills (skill_id)
) ENGINE=InnoDB COMMENT='职位技能要求';

CREATE TABLE candidates (
  candidate_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  candidate_no VARCHAR(30) NOT NULL,
  full_name VARCHAR(80) NOT NULL,
  gender ENUM('female','male','non_binary','undisclosed') NOT NULL DEFAULT 'undisclosed',
  phone VARCHAR(30) NOT NULL,
  email VARCHAR(160) NOT NULL,
  current_city VARCHAR(60) NOT NULL,
  current_company VARCHAR(120) NULL,
  current_title VARCHAR(120) NULL,
  years_of_experience DECIMAL(4,1) NOT NULL DEFAULT 0,
  highest_education ENUM('high_school','associate','bachelor','master','doctorate') NOT NULL,
  source_id BIGINT UNSIGNED NOT NULL,
  consent_to_store BOOLEAN NOT NULL DEFAULT TRUE,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (candidate_id),
  UNIQUE KEY uk_candidates_no (candidate_no),
  UNIQUE KEY uk_candidates_email (email),
  KEY idx_candidates_source (source_id),
  KEY idx_candidates_city_experience (current_city, years_of_experience),
  CONSTRAINT fk_candidates_source FOREIGN KEY (source_id)
    REFERENCES candidate_sources (source_id)
) ENGINE=InnoDB COMMENT='候选人主档（全部为虚构数据）';

CREATE TABLE candidate_educations (
  education_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  candidate_id BIGINT UNSIGNED NOT NULL,
  school_name VARCHAR(160) NOT NULL,
  degree_level ENUM('associate','bachelor','master','doctorate') NOT NULL,
  major VARCHAR(120) NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NULL,
  PRIMARY KEY (education_id),
  KEY idx_candidate_educations_candidate (candidate_id),
  CONSTRAINT fk_candidate_educations_candidate FOREIGN KEY (candidate_id)
    REFERENCES candidates (candidate_id) ON DELETE CASCADE
) ENGINE=InnoDB COMMENT='候选人教育经历';

CREATE TABLE candidate_experiences (
  experience_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  candidate_id BIGINT UNSIGNED NOT NULL,
  company_name VARCHAR(160) NOT NULL,
  job_title VARCHAR(120) NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NULL,
  is_current BOOLEAN NOT NULL DEFAULT FALSE,
  summary VARCHAR(500) NOT NULL,
  PRIMARY KEY (experience_id),
  KEY idx_candidate_experiences_candidate (candidate_id, start_date),
  CONSTRAINT fk_candidate_experiences_candidate FOREIGN KEY (candidate_id)
    REFERENCES candidates (candidate_id) ON DELETE CASCADE
) ENGINE=InnoDB COMMENT='候选人工作经历';

CREATE TABLE applications (
  application_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  application_no VARCHAR(30) NOT NULL,
  candidate_id BIGINT UNSIGNED NOT NULL,
  requisition_id BIGINT UNSIGNED NOT NULL,
  recruiter_id BIGINT UNSIGNED NOT NULL,
  applied_at DATETIME NOT NULL,
  current_stage ENUM('applied','resume_review','phone_screen','interview','final_interview','offer','hired','closed') NOT NULL,
  status ENUM('active','rejected','withdrawn','offer_extended','offer_accepted','offer_declined','hired') NOT NULL,
  expected_salary DECIMAL(12,2) NULL,
  available_date DATE NULL,
  rejection_reason VARCHAR(255) NULL,
  withdrawn_reason VARCHAR(255) NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (application_id),
  UNIQUE KEY uk_applications_no (application_no),
  UNIQUE KEY uk_applications_candidate_requisition (candidate_id, requisition_id),
  KEY idx_applications_requisition_status (requisition_id, status),
  KEY idx_applications_stage_applied (current_stage, applied_at),
  KEY idx_applications_recruiter (recruiter_id),
  CONSTRAINT fk_applications_candidate FOREIGN KEY (candidate_id)
    REFERENCES candidates (candidate_id),
  CONSTRAINT fk_applications_requisition FOREIGN KEY (requisition_id)
    REFERENCES job_requisitions (requisition_id),
  CONSTRAINT fk_applications_recruiter FOREIGN KEY (recruiter_id)
    REFERENCES employees (employee_id)
) ENGINE=InnoDB COMMENT='候选人的职位申请';

CREATE TABLE application_status_history (
  history_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  application_id BIGINT UNSIGNED NOT NULL,
  from_stage VARCHAR(30) NULL,
  to_stage VARCHAR(30) NOT NULL,
  from_status VARCHAR(30) NULL,
  to_status VARCHAR(30) NOT NULL,
  changed_by_employee_id BIGINT UNSIGNED NOT NULL,
  changed_at DATETIME NOT NULL,
  change_reason VARCHAR(255) NULL,
  PRIMARY KEY (history_id),
  KEY idx_application_history_application_time (application_id, changed_at),
  CONSTRAINT fk_application_history_application FOREIGN KEY (application_id)
    REFERENCES applications (application_id) ON DELETE CASCADE,
  CONSTRAINT fk_application_history_employee FOREIGN KEY (changed_by_employee_id)
    REFERENCES employees (employee_id)
) ENGINE=InnoDB COMMENT='申请阶段与状态变更历史';

CREATE TABLE interview_types (
  interview_type_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  type_code VARCHAR(30) NOT NULL,
  type_name VARCHAR(80) NOT NULL,
  default_duration_minutes SMALLINT UNSIGNED NOT NULL,
  sequence_hint TINYINT UNSIGNED NOT NULL,
  PRIMARY KEY (interview_type_id),
  UNIQUE KEY uk_interview_types_code (type_code)
) ENGINE=InnoDB COMMENT='面试类型';

CREATE TABLE interview_sessions (
  interview_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  interview_no VARCHAR(30) NOT NULL,
  application_id BIGINT UNSIGNED NOT NULL,
  interview_type_id BIGINT UNSIGNED NOT NULL,
  round_number TINYINT UNSIGNED NOT NULL,
  scheduled_start DATETIME NOT NULL,
  scheduled_end DATETIME NOT NULL,
  timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Shanghai',
  mode ENUM('onsite','video','phone') NOT NULL,
  office_id BIGINT UNSIGNED NULL,
  meeting_room VARCHAR(80) NULL,
  meeting_link VARCHAR(255) NULL,
  coordinator_id BIGINT UNSIGNED NOT NULL,
  status ENUM('scheduled','confirmed','completed','cancelled','no_show','reschedule_requested') NOT NULL,
  candidate_confirmed_at DATETIME NULL,
  completed_at DATETIME NULL,
  cancellation_reason VARCHAR(255) NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (interview_id),
  UNIQUE KEY uk_interview_sessions_no (interview_no),
  UNIQUE KEY uk_interview_sessions_round (application_id, round_number),
  KEY idx_interviews_schedule_status (scheduled_start, status),
  KEY idx_interviews_coordinator (coordinator_id),
  CONSTRAINT chk_interview_time CHECK (scheduled_end > scheduled_start),
  CONSTRAINT fk_interviews_application FOREIGN KEY (application_id)
    REFERENCES applications (application_id) ON DELETE CASCADE,
  CONSTRAINT fk_interviews_type FOREIGN KEY (interview_type_id)
    REFERENCES interview_types (interview_type_id),
  CONSTRAINT fk_interviews_office FOREIGN KEY (office_id)
    REFERENCES offices (office_id),
  CONSTRAINT fk_interviews_coordinator FOREIGN KEY (coordinator_id)
    REFERENCES employees (employee_id)
) ENGINE=InnoDB COMMENT='具体面试日程';

CREATE TABLE interview_panel_members (
  interview_id BIGINT UNSIGNED NOT NULL,
  employee_id BIGINT UNSIGNED NOT NULL,
  panel_role ENUM('lead','interviewer','observer','bar_raiser') NOT NULL,
  attendance_status ENUM('pending','accepted','declined','attended','absent') NOT NULL DEFAULT 'pending',
  assigned_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (interview_id, employee_id),
  KEY idx_panel_members_employee (employee_id, attendance_status),
  CONSTRAINT fk_panel_members_interview FOREIGN KEY (interview_id)
    REFERENCES interview_sessions (interview_id) ON DELETE CASCADE,
  CONSTRAINT fk_panel_members_employee FOREIGN KEY (employee_id)
    REFERENCES employees (employee_id)
) ENGINE=InnoDB COMMENT='面试小组成员';

CREATE TABLE competencies (
  competency_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  competency_code VARCHAR(30) NOT NULL,
  competency_name VARCHAR(100) NOT NULL,
  category ENUM('general','technical','leadership','culture') NOT NULL,
  description VARCHAR(500) NOT NULL,
  PRIMARY KEY (competency_id),
  UNIQUE KEY uk_competencies_code (competency_code)
) ENGINE=InnoDB COMMENT='面试评分能力项';

CREATE TABLE interview_feedback (
  feedback_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  interview_id BIGINT UNSIGNED NOT NULL,
  interviewer_id BIGINT UNSIGNED NOT NULL,
  overall_rating TINYINT UNSIGNED NOT NULL,
  recommendation ENUM('strong_no','no','mixed','yes','strong_yes') NOT NULL,
  strengths TEXT NOT NULL,
  concerns TEXT NOT NULL,
  evidence_notes TEXT NOT NULL,
  submitted_at DATETIME NOT NULL,
  is_final BOOLEAN NOT NULL DEFAULT TRUE,
  PRIMARY KEY (feedback_id),
  UNIQUE KEY uk_feedback_interview_interviewer (interview_id, interviewer_id),
  KEY idx_feedback_recommendation (recommendation, submitted_at),
  CONSTRAINT chk_feedback_rating CHECK (overall_rating BETWEEN 1 AND 5),
  CONSTRAINT fk_feedback_interview FOREIGN KEY (interview_id)
    REFERENCES interview_sessions (interview_id) ON DELETE CASCADE,
  CONSTRAINT fk_feedback_interviewer FOREIGN KEY (interviewer_id)
    REFERENCES employees (employee_id)
) ENGINE=InnoDB COMMENT='面试官综合反馈';

CREATE TABLE feedback_competency_scores (
  feedback_id BIGINT UNSIGNED NOT NULL,
  competency_id BIGINT UNSIGNED NOT NULL,
  score TINYINT UNSIGNED NOT NULL,
  comment VARCHAR(500) NOT NULL,
  PRIMARY KEY (feedback_id, competency_id),
  CONSTRAINT chk_competency_score CHECK (score BETWEEN 1 AND 5),
  CONSTRAINT fk_competency_scores_feedback FOREIGN KEY (feedback_id)
    REFERENCES interview_feedback (feedback_id) ON DELETE CASCADE,
  CONSTRAINT fk_competency_scores_competency FOREIGN KEY (competency_id)
    REFERENCES competencies (competency_id)
) ENGINE=InnoDB COMMENT='各能力项的细分评分';

CREATE TABLE offers (
  offer_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  offer_no VARCHAR(30) NOT NULL,
  application_id BIGINT UNSIGNED NOT NULL,
  offered_title VARCHAR(120) NOT NULL,
  job_level VARCHAR(30) NOT NULL,
  base_salary DECIMAL(12,2) NOT NULL,
  annual_bonus_target DECIMAL(5,2) NOT NULL DEFAULT 0,
  sign_on_bonus DECIMAL(12,2) NOT NULL DEFAULT 0,
  equity_value DECIMAL(12,2) NOT NULL DEFAULT 0,
  currency CHAR(3) NOT NULL DEFAULT 'CNY',
  proposed_start_date DATE NOT NULL,
  expires_at DATETIME NOT NULL,
  status ENUM('draft','pending_approval','approved','sent','accepted','declined','expired','withdrawn') NOT NULL,
  sent_at DATETIME NULL,
  responded_at DATETIME NULL,
  decline_reason VARCHAR(255) NULL,
  created_by_employee_id BIGINT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (offer_id),
  UNIQUE KEY uk_offers_no (offer_no),
  UNIQUE KEY uk_offers_application (application_id),
  KEY idx_offers_status_expiry (status, expires_at),
  CONSTRAINT fk_offers_application FOREIGN KEY (application_id)
    REFERENCES applications (application_id),
  CONSTRAINT fk_offers_creator FOREIGN KEY (created_by_employee_id)
    REFERENCES employees (employee_id)
) ENGINE=InnoDB COMMENT='录用意向与候选人响应';

CREATE TABLE offer_approvals (
  approval_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  offer_id BIGINT UNSIGNED NOT NULL,
  approval_step TINYINT UNSIGNED NOT NULL,
  approver_id BIGINT UNSIGNED NOT NULL,
  decision ENUM('pending','approved','rejected','skipped') NOT NULL,
  decided_at DATETIME NULL,
  comment VARCHAR(500) NULL,
  PRIMARY KEY (approval_id),
  UNIQUE KEY uk_offer_approval_step (offer_id, approval_step),
  KEY idx_offer_approvals_approver (approver_id, decision),
  CONSTRAINT fk_offer_approvals_offer FOREIGN KEY (offer_id)
    REFERENCES offers (offer_id) ON DELETE CASCADE,
  CONSTRAINT fk_offer_approvals_approver FOREIGN KEY (approver_id)
    REFERENCES employees (employee_id)
) ENGINE=InnoDB COMMENT='Offer 审批链';

CREATE VIEW v_application_overview AS
SELECT
  a.application_id,
  a.application_no,
  c.candidate_no,
  c.full_name AS candidate_name,
  c.current_city,
  jr.requisition_no,
  jr.job_title,
  d.department_name,
  r.full_name AS recruiter_name,
  a.applied_at,
  a.current_stage,
  a.status,
  COUNT(DISTINCT i.interview_id) AS interview_count,
  ROUND(AVG(f.overall_rating), 2) AS average_rating
FROM applications a
JOIN candidates c ON c.candidate_id = a.candidate_id
JOIN job_requisitions jr ON jr.requisition_id = a.requisition_id
JOIN departments d ON d.department_id = jr.department_id
JOIN employees r ON r.employee_id = a.recruiter_id
LEFT JOIN interview_sessions i ON i.application_id = a.application_id
LEFT JOIN interview_feedback f ON f.interview_id = i.interview_id
GROUP BY a.application_id, a.application_no, c.candidate_no, c.full_name,
  c.current_city, jr.requisition_no, jr.job_title, d.department_name,
  r.full_name, a.applied_at, a.current_stage, a.status;

CREATE VIEW v_interview_schedule AS
SELECT
  i.interview_id,
  i.interview_no,
  i.scheduled_start,
  i.scheduled_end,
  i.mode,
  i.status,
  it.type_name AS interview_type,
  c.full_name AS candidate_name,
  jr.job_title,
  o.office_name,
  i.meeting_room,
  COUNT(pm.employee_id) AS panel_size
FROM interview_sessions i
JOIN interview_types it ON it.interview_type_id = i.interview_type_id
JOIN applications a ON a.application_id = i.application_id
JOIN candidates c ON c.candidate_id = a.candidate_id
JOIN job_requisitions jr ON jr.requisition_id = a.requisition_id
LEFT JOIN offices o ON o.office_id = i.office_id
LEFT JOIN interview_panel_members pm ON pm.interview_id = i.interview_id
GROUP BY i.interview_id, i.interview_no, i.scheduled_start, i.scheduled_end,
  i.mode, i.status, it.type_name, c.full_name, jr.job_title,
  o.office_name, i.meeting_room;

CREATE VIEW v_recruiting_funnel AS
SELECT
  jr.requisition_id,
  jr.requisition_no,
  jr.job_title,
  d.department_name,
  COUNT(a.application_id) AS total_applications,
  SUM(a.current_stage IN ('phone_screen','interview','final_interview','offer','hired')) AS screened_count,
  SUM(a.current_stage IN ('interview','final_interview','offer','hired')) AS interviewed_count,
  SUM(a.current_stage IN ('offer','hired')) AS offer_count,
  SUM(a.status = 'hired') AS hired_count,
  ROUND(100 * SUM(a.status = 'hired') / NULLIF(COUNT(a.application_id), 0), 2) AS hire_rate_pct
FROM job_requisitions jr
JOIN departments d ON d.department_id = jr.department_id
LEFT JOIN applications a ON a.requisition_id = jr.requisition_id
GROUP BY jr.requisition_id, jr.requisition_no, jr.job_title, d.department_name;

CREATE VIEW v_interviewer_workload AS
SELECT
  e.employee_id,
  e.full_name AS interviewer_name,
  d.department_name,
  COUNT(pm.interview_id) AS assigned_interviews,
  SUM(pm.attendance_status = 'attended') AS attended_interviews,
  COUNT(f.feedback_id) AS submitted_feedback,
  ROUND(AVG(f.overall_rating), 2) AS average_rating_given
FROM employees e
JOIN departments d ON d.department_id = e.department_id
LEFT JOIN interview_panel_members pm ON pm.employee_id = e.employee_id
LEFT JOIN interview_feedback f
  ON f.interview_id = pm.interview_id AND f.interviewer_id = e.employee_id
GROUP BY e.employee_id, e.full_name, d.department_name;

REVOKE ALL PRIVILEGES ON interview_db.* FROM 'interview_reader'@'%';
GRANT SELECT ON interview_db.* TO 'interview_reader'@'%';
FLUSH PRIVILEGES;
