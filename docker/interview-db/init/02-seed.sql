-- Deterministic synthetic seed data, snapshot date: 2026-07-22.
-- It is intended for local demos, SQL practice, and integration testing only.

USE interview_db;
SET NAMES utf8mb4;
SET time_zone = '+08:00';

INSERT INTO offices
  (office_id, office_code, office_name, city, address, timezone)
VALUES
  (1, 'SH-HQ', '上海总部', '上海', '浦东新区张江路 88 号', 'Asia/Shanghai'),
  (2, 'BJ-RD', '北京研发中心', '北京', '海淀区知春路 66 号', 'Asia/Shanghai'),
  (3, 'SZ-TECH', '深圳技术中心', '深圳', '南山区科技园科苑路 18 号', 'Asia/Shanghai'),
  (4, 'HZ-PROD', '杭州产品中心', '杭州', '余杭区文一西路 998 号', 'Asia/Shanghai'),
  (5, 'CD-CS', '成都客户中心', '成都', '高新区天府五街 168 号', 'Asia/Shanghai');

INSERT INTO departments
  (department_id, parent_department_id, department_code, department_name, cost_center)
VALUES
  (1, NULL, 'ENG', '技术中心', 'CC1000'),
  (2, NULL, 'PROD', '产品中心', 'CC2000'),
  (3, NULL, 'DATA', '数据智能中心', 'CC3000'),
  (4, NULL, 'GTM', '市场与商业化中心', 'CC4000'),
  (5, NULL, 'CORP', '企业职能中心', 'CC5000'),
  (6, 1, 'ENG-BE', '后端研发部', 'CC1100'),
  (7, 1, 'ENG-FE', '前端研发部', 'CC1200'),
  (8, 1, 'ENG-QA', '质量工程部', 'CC1300'),
  (9, 1, 'ENG-SRE', '平台与稳定性部', 'CC1400'),
  (10, 2, 'PROD-PLAT', '平台产品部', 'CC2100'),
  (11, 2, 'PROD-UX', '用户体验设计部', 'CC2200'),
  (12, 3, 'DATA-ALG', '算法工程部', 'CC3100'),
  (13, 3, 'DATA-BI', '商业分析部', 'CC3200'),
  (14, 4, 'GTM-SALES', '企业销售部', 'CC4100'),
  (15, 4, 'GTM-MKT', '品牌市场部', 'CC4200'),
  (16, 5, 'CORP-HR', '人力资源部', 'CC5100'),
  (17, 5, 'CORP-FIN', '财务部', 'CC5200');

INSERT INTO employees
  (employee_id, employee_no, full_name, email, department_id, office_id, job_title, employee_role, hire_date)
VALUES
  (1, 'E1001', '周欣怡', 'xinyi.zhou@example.test', 16, 1, '高级招聘经理', 'recruiter', '2019-03-11'),
  (2, 'E1002', '陈嘉宁', 'jianing.chen@example.test', 16, 2, '技术招聘专家', 'recruiter', '2020-06-15'),
  (3, 'E1003', '徐思远', 'siyuan.xu@example.test', 16, 3, '高级招聘顾问', 'recruiter', '2021-02-22'),
  (4, 'E1004', '梁雨桐', 'yutong.liang@example.test', 16, 4, '产品招聘专家', 'recruiter', '2021-09-06'),
  (5, 'E1005', '沈文博', 'wenbo.shen@example.test', 16, 1, '商业招聘专家', 'recruiter', '2022-04-18'),
  (6, 'E1006', '蒋若琳', 'ruolin.jiang@example.test', 16, 5, '招聘运营专员', 'recruiter', '2023-07-10'),
  (7, 'E1101', '顾承泽', 'chengze.gu@example.test', 6, 1, '后端研发总监', 'hiring_manager', '2017-05-08'),
  (8, 'E1102', '韩子墨', 'zimo.han@example.test', 7, 2, '前端研发经理', 'hiring_manager', '2018-08-20'),
  (9, 'E1103', '赵清扬', 'qingyang.zhao@example.test', 8, 3, '质量工程经理', 'hiring_manager', '2018-11-12'),
  (10, 'E1104', '林致远', 'zhiyuan.lin@example.test', 9, 2, 'SRE 负责人', 'hiring_manager', '2019-01-21'),
  (11, 'E1201', '陆安然', 'anran.lu@example.test', 10, 4, '平台产品总监', 'hiring_manager', '2018-03-05'),
  (12, 'E1202', '宋知夏', 'zhixia.song@example.test', 11, 4, '设计负责人', 'hiring_manager', '2019-07-01'),
  (13, 'E1301', '谢明哲', 'mingzhe.xie@example.test', 12, 2, '算法总监', 'hiring_manager', '2017-10-16'),
  (14, 'E1302', '方可欣', 'kexin.fang@example.test', 13, 1, '数据分析经理', 'hiring_manager', '2020-01-06'),
  (15, 'E1401', '袁浩然', 'haoran.yuan@example.test', 14, 1, '销售总监', 'hiring_manager', '2018-05-14'),
  (16, 'E1402', '唐婉清', 'wanqing.tang@example.test', 15, 1, '市场总监', 'hiring_manager', '2019-12-02'),
  (17, 'E2101', '程慕白', 'mubai.cheng@example.test', 6, 1, '资深 Go 工程师', 'interviewer', '2020-04-13'),
  (18, 'E2102', '苏景行', 'jingxing.su@example.test', 6, 2, '资深 Java 工程师', 'interviewer', '2019-06-24'),
  (19, 'E2103', '罗一诺', 'yinuo.luo@example.test', 6, 3, '架构师', 'interviewer', '2018-09-17'),
  (20, 'E2201', '魏舒然', 'shuran.wei@example.test', 7, 2, '资深前端工程师', 'interviewer', '2020-08-03'),
  (21, 'E2202', '叶书瑶', 'shuyao.ye@example.test', 7, 4, '前端技术专家', 'interviewer', '2019-02-11'),
  (22, 'E2301', '邹逸辰', 'yichen.zou@example.test', 8, 3, '高级测试开发工程师', 'interviewer', '2021-05-17'),
  (23, 'E2302', '秦乐之', 'lezhi.qin@example.test', 9, 2, '资深 SRE 工程师', 'interviewer', '2020-10-26'),
  (24, 'E2401', '孟星河', 'xinghe.meng@example.test', 10, 4, '高级产品经理', 'interviewer', '2021-01-18'),
  (25, 'E2402', '夏语冰', 'yubing.xia@example.test', 11, 4, '资深交互设计师', 'interviewer', '2020-03-09'),
  (26, 'E2501', '白敬言', 'jingyan.bai@example.test', 12, 2, '推荐算法专家', 'interviewer', '2018-07-23'),
  (27, 'E2502', '任诗涵', 'shihan.ren@example.test', 12, 1, '机器学习平台专家', 'interviewer', '2019-11-04'),
  (28, 'E2503', '金嘉树', 'jiashu.jin@example.test', 13, 1, '资深数据分析师', 'interviewer', '2021-06-28'),
  (29, 'E2601', '段思齐', 'siqi.duan@example.test', 14, 1, '大客户销售经理', 'interviewer', '2020-12-14'),
  (30, 'E2602', '何曼文', 'manwen.he@example.test', 15, 1, '品牌策略经理', 'interviewer', '2021-03-22'),
  (31, 'E3001', '吴海峰', 'haifeng.wu@example.test', 5, 1, '人力资源副总裁', 'approver', '2016-04-06'),
  (32, 'E3002', '郑雅琴', 'yaqin.zheng@example.test', 17, 1, '财务总监', 'approver', '2017-09-25'),
  (33, 'E3003', '冯天宇', 'tianyu.feng@example.test', 16, 1, 'HRBP 负责人', 'hrbp', '2019-05-27'),
  (34, 'E3004', '曹静姝', 'jingshu.cao@example.test', 16, 2, '技术 HRBP', 'hrbp', '2020-07-13'),
  (35, 'E3005', '孔德馨', 'dexin.kong@example.test', 16, 4, '产品 HRBP', 'hrbp', '2021-08-30');

INSERT INTO candidate_sources
  (source_id, source_code, source_name, source_category, cost_per_hire)
VALUES
  (1, 'CAREER', '公司招聘官网', 'career_site', 350.00),
  (2, 'REFERRAL', '内部员工推荐', 'referral', 5000.00),
  (3, 'BOSS', 'BOSS 直聘', 'job_board', 6800.00),
  (4, 'LIEPIN', '猎聘', 'job_board', 9200.00),
  (5, 'LINKEDIN', 'LinkedIn', 'social', 7600.00),
  (6, 'CAMPUS', '校园招聘', 'campus', 3100.00),
  (7, 'AGENCY-A', '启航猎头', 'agency', 28000.00),
  (8, 'AGENCY-B', '远见人才顾问', 'agency', 32000.00),
  (9, 'TECH-MEETUP', '技术社区活动', 'event', 2200.00),
  (10, 'WECHAT', '微信公众号', 'social', 900.00);

INSERT INTO skills (skill_id, skill_name, skill_category) VALUES
  (1, 'Go', 'engineering'), (2, 'Java', 'engineering'),
  (3, 'Python', 'engineering'), (4, 'TypeScript', 'engineering'),
  (5, 'Vue.js', 'engineering'), (6, 'React', 'engineering'),
  (7, 'MySQL', 'engineering'), (8, 'Redis', 'engineering'),
  (9, 'Kafka', 'engineering'), (10, 'Kubernetes', 'engineering'),
  (11, '微服务架构', 'engineering'), (12, '自动化测试', 'engineering'),
  (13, '产品规划', 'product'), (14, '用户研究', 'product'),
  (15, 'Figma', 'design'), (16, '交互设计', 'design'),
  (17, '机器学习', 'data'), (18, '深度学习', 'data'),
  (19, 'SQL', 'data'), (20, '数据可视化', 'data'),
  (21, 'A/B 测试', 'data'), (22, '企业销售', 'business'),
  (23, '品牌营销', 'business'), (24, '商务谈判', 'business'),
  (25, '英语', 'language'), (26, '项目管理', 'leadership'),
  (27, '团队管理', 'leadership'), (28, '系统设计', 'engineering'),
  (29, 'Linux', 'engineering'), (30, '云原生', 'engineering');

INSERT INTO interview_types
  (interview_type_id, type_code, type_name, default_duration_minutes, sequence_hint)
VALUES
  (1, 'HR_SCREEN', 'HR 初筛', 30, 1),
  (2, 'TECHNICAL', '专业技术面试', 60, 2),
  (3, 'CASE_STUDY', '案例分析', 75, 2),
  (4, 'SYSTEM_DESIGN', '系统设计面试', 60, 3),
  (5, 'MANAGER', '直属经理面试', 45, 3),
  (6, 'CULTURE', '价值观面试', 45, 4),
  (7, 'FINAL', '终面', 45, 5);

INSERT INTO competencies
  (competency_id, competency_code, competency_name, category, description)
VALUES
  (1, 'PROBLEM_SOLVING', '问题分析与解决', 'general', '能否拆解复杂问题并形成可验证的解决方案'),
  (2, 'COMMUNICATION', '沟通表达', 'general', '表达是否清晰、结构化并能有效倾听'),
  (3, 'DOMAIN_DEPTH', '专业深度', 'technical', '核心专业知识的广度、深度与准确性'),
  (4, 'SYSTEM_THINKING', '系统性思维', 'technical', '能否识别边界、依赖、容量及风险'),
  (5, 'EXECUTION', '结果交付', 'general', '能否在约束下推动工作并稳定交付'),
  (6, 'COLLABORATION', '团队协作', 'culture', '是否愿意共享信息、处理分歧并支持团队'),
  (7, 'CUSTOMER_FOCUS', '客户意识', 'culture', '是否理解客户价值并据此做出取舍'),
  (8, 'LEARNING', '学习敏捷性', 'culture', '面对新领域时获取反馈并快速迭代的能力'),
  (9, 'LEADERSHIP', '领导力', 'leadership', '设定方向、影响他人及承担责任的能力'),
  (10, 'VALUES', '价值观契合', 'culture', '行为方式与公司价值观的匹配程度');

DELIMITER $$

CREATE PROCEDURE seed_interview_db()
BEGIN
  DECLARE i INT DEFAULT 1;
  DECLARE r INT DEFAULT 1;
  DECLARE v_department BIGINT UNSIGNED;
  DECLARE v_office BIGINT UNSIGNED;
  DECLARE v_manager BIGINT UNSIGNED;
  DECLARE v_recruiter BIGINT UNSIGNED;
  DECLARE v_title VARCHAR(120);
  DECLARE v_level VARCHAR(30);
  DECLARE v_min_salary DECIMAL(12,2);
  DECLARE v_candidate BIGINT UNSIGNED;
  DECLARE v_requisition BIGINT UNSIGNED;
  DECLARE v_stage VARCHAR(30);
  DECLARE v_status VARCHAR(30);
  DECLARE v_applied_at DATETIME;
  DECLARE v_rounds INT;
  DECLARE v_interview_status VARCHAR(30);
  DECLARE v_interview_start DATETIME;
  DECLARE v_mode VARCHAR(20);
  DECLARE v_interview_id BIGINT UNSIGNED;
  DECLARE v_interviewer_1 BIGINT UNSIGNED;
  DECLARE v_interviewer_2 BIGINT UNSIGNED;

  WHILE i <= 24 DO
    SET v_department = ELT(((i - 1) MOD 10) + 1, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15);
    SET v_office = ELT(((i - 1) MOD 5) + 1, 1, 2, 3, 4, 5);
    SET v_manager = CASE v_department
      WHEN 6 THEN 7 WHEN 7 THEN 8 WHEN 8 THEN 9 WHEN 9 THEN 10
      WHEN 10 THEN 11 WHEN 11 THEN 12 WHEN 12 THEN 13 WHEN 13 THEN 14
      WHEN 14 THEN 15 ELSE 16 END;
    SET v_recruiter = ((i - 1) MOD 6) + 1;
    SET v_title = ELT(((i - 1) MOD 24) + 1,
      'Go 后端工程师', '高级 Java 工程师', '前端开发工程师', '测试开发工程师',
      'SRE 工程师', '平台产品经理', '高级交互设计师', '机器学习工程师',
      '数据分析师', '大客户销售经理', '品牌营销经理', '云原生架构师',
      '资深后端工程师', '前端技术专家', '自动化测试专家', '运维平台工程师',
      '商业产品经理', '视觉设计师', '推荐算法工程师', '商业分析师',
      '销售解决方案顾问', '内容运营经理', '数据平台工程师', '研发经理');
    SET v_level = ELT(((i - 1) MOD 6) + 1, 'mid', 'senior', 'mid', 'senior', 'lead', 'manager');
    SET v_min_salary = 180000 + (((i - 1) MOD 8) * 30000);

    INSERT INTO job_requisitions
      (requisition_no, job_title, department_id, office_id, hiring_manager_id,
       recruiter_id, employment_type, workplace_type, job_level, headcount,
       filled_count, min_salary, max_salary, status, opened_at,
       target_close_date, closed_at, description)
    VALUES
      (CONCAT('REQ-2026-', LPAD(i, 3, '0')), v_title, v_department, v_office,
       v_manager, v_recruiter,
       IF(i IN (6, 18), 'contract', IF(i IN (11, 22), 'intern', 'full_time')),
       ELT(((i - 1) MOD 3) + 1, 'hybrid', 'onsite', 'remote'), v_level,
       IF(i MOD 7 = 0, 3, IF(i MOD 4 = 0, 2, 1)),
       IF(i BETWEEN 19 AND 21, 1, 0), v_min_salary, v_min_salary + 120000,
       CASE WHEN i <= 18 THEN 'open' WHEN i <= 21 THEN 'closed'
            WHEN i = 22 THEN 'on_hold' WHEN i = 23 THEN 'closed' ELSE 'cancelled' END,
       DATE_ADD('2026-01-06', INTERVAL (i * 4) DAY),
       DATE_ADD('2026-01-06', INTERVAL (i * 4 + 75) DAY),
       IF(i IN (19,20,21,23), DATE_ADD('2026-01-06', INTERVAL (i * 4 + 58) DAY), NULL),
       CONCAT('负责', v_title, '相关工作；关注业务结果、跨团队协作和可持续交付。'));
    SET i = i + 1;
  END WHILE;

  INSERT INTO requisition_skills (requisition_id, skill_id, importance, minimum_years)
  SELECT requisition_id, ((requisition_id - 1) MOD 30) + 1, 'required', 3.0
  FROM job_requisitions;
  INSERT INTO requisition_skills (requisition_id, skill_id, importance, minimum_years)
  SELECT requisition_id, ((requisition_id + 6) MOD 30) + 1, 'required', 2.0
  FROM job_requisitions;
  INSERT INTO requisition_skills (requisition_id, skill_id, importance, minimum_years)
  SELECT requisition_id, ((requisition_id + 13) MOD 30) + 1, 'preferred', 1.0
  FROM job_requisitions;

  SET i = 1;
  WHILE i <= 160 DO
    INSERT INTO candidates
      (candidate_no, full_name, gender, phone, email, current_city, current_company,
       current_title, years_of_experience, highest_education, source_id, created_at)
    VALUES
      (CONCAT('CAN-2026-', LPAD(i, 4, '0')),
       CONCAT(
         ELT(((i - 1) MOD 20) + 1, '王','李','张','刘','陈','杨','黄','赵','周','吴','徐','孙','胡','朱','高','林','何','郭','马','罗'),
         ELT(((i * 7 - 1) MOD 24) + 1, '子涵','宇轩','雨桐','思远','若曦','浩然','嘉宁','欣怡','明哲','可欣','景行','舒然','安琪','致远','诗涵','一诺','星河','书瑶','慕白','清扬','婉清','逸辰','乐之','语冰')),
       ELT(((i - 1) MOD 4) + 1, 'female', 'male', 'undisclosed', 'male'),
       CONCAT('1', ELT(((i - 1) MOD 5) + 1, '38','39','86','88','89'), LPAD(10000000 + i, 8, '0')),
       CONCAT('candidate', LPAD(i, 4, '0'), '@example.test'),
       ELT(((i - 1) MOD 10) + 1, '上海','北京','深圳','杭州','成都','南京','苏州','武汉','西安','广州'),
       ELT(((i * 3 - 1) MOD 15) + 1, '星云科技','远景网络','青禾数据','极光软件','海岳信息','数澜智能','云桥科技','明日互动','北辰系统','新岸咨询','致远电商','凌云汽车','百川金融','创想互联','启明数字'),
       ELT(((i * 5 - 1) MOD 16) + 1, '后端工程师','高级软件工程师','前端工程师','测试开发工程师','SRE 工程师','产品经理','交互设计师','算法工程师','数据分析师','客户经理','市场经理','技术负责人','业务分析师','解决方案顾问','数据工程师','研发经理'),
       ((i * 7) MOD 13) + 0.5,
       ELT(((i - 1) MOD 8) + 1, 'bachelor','bachelor','master','bachelor','master','associate','bachelor','doctorate'),
       ((i - 1) MOD 10) + 1,
       DATE_ADD('2026-01-02 09:00:00', INTERVAL (i * 19) HOUR));

    INSERT INTO candidate_educations
      (candidate_id, school_name, degree_level, major, start_date, end_date)
    VALUES
      (i,
       ELT(((i - 1) MOD 16) + 1, '华东理工大学','北京交通大学','深圳大学','浙江工业大学','电子科技大学','南京理工大学','武汉理工大学','西北大学','暨南大学','苏州大学','北京工业大学','上海大学','杭州电子科技大学','西安电子科技大学','华南理工大学','重庆大学'),
       ELT(((i - 1) MOD 8) + 1, 'bachelor','bachelor','master','bachelor','master','associate','bachelor','doctorate'),
       ELT(((i * 3 - 1) MOD 10) + 1, '计算机科学与技术','软件工程','信息管理','电子信息工程','统计学','工业设计','工商管理','市场营销','数学与应用数学','自动化'),
       MAKEDATE(2006 + (i MOD 11), 244),
       MAKEDATE(2010 + (i MOD 11), 181));

    INSERT INTO candidate_experiences
      (candidate_id, company_name, job_title, start_date, end_date, is_current, summary)
    VALUES
      (i,
       ELT(((i * 3 - 1) MOD 15) + 1, '星云科技','远景网络','青禾数据','极光软件','海岳信息','数澜智能','云桥科技','明日互动','北辰系统','新岸咨询','致远电商','凌云汽车','百川金融','创想互联','启明数字'),
       ELT(((i * 5 - 1) MOD 16) + 1, '后端工程师','高级软件工程师','前端工程师','测试开发工程师','SRE 工程师','产品经理','交互设计师','算法工程师','数据分析师','客户经理','市场经理','技术负责人','业务分析师','解决方案顾问','数据工程师','研发经理'),
       DATE_SUB('2026-06-30', INTERVAL (12 + (i MOD 48)) MONTH), NULL, TRUE,
       ELT(((i - 1) MOD 6) + 1,
         '负责核心交易服务的迭代、稳定性治理与性能优化。',
         '负责企业客户需求分析、方案设计及跨团队项目交付。',
         '参与数据平台建设，完善指标体系和自助分析能力。',
         '负责产品规划、用户调研、版本管理和效果复盘。',
         '建设自动化测试与发布流程，降低线上缺陷率。',
         '负责算法建模、离线评估、在线实验和持续优化。'));

    IF i MOD 3 <> 0 THEN
      INSERT INTO candidate_experiences
        (candidate_id, company_name, job_title, start_date, end_date, is_current, summary)
      VALUES
        (i,
         ELT(((i * 7 - 1) MOD 12) + 1, '原点科技','云帆软件','天际网络','博思数据','晨星咨询','新城互联','知微信息','开源动力','恒川科技','蓝图数字','领航系统','同创软件'),
         ELT(((i * 2 - 1) MOD 8) + 1, '软件工程师','初级产品经理','数据专员','测试工程师','实施顾问','运营专员','算法工程师','销售顾问'),
         DATE_SUB('2026-06-30', INTERVAL (55 + (i MOD 40)) MONTH),
         DATE_SUB('2026-06-30', INTERVAL (14 + (i MOD 38)) MONTH), FALSE,
         '承担模块交付、问题跟踪、文档沉淀及跨职能协作。');
    END IF;
    SET i = i + 1;
  END WHILE;

  SET i = 1;
  WHILE i <= 220 DO
    SET v_candidate = IF(i <= 160, i, i - 160);
    SET v_requisition = IF(i <= 160, ((v_candidate - 1) MOD 24) + 1, ((v_candidate + 7) MOD 24) + 1);
    SET v_recruiter = ((v_requisition - 1) MOD 6) + 1;
    SET v_applied_at = DATE_ADD('2026-01-05 09:30:00', INTERVAL (i * 17) HOUR);
    SET v_status = CASE (i MOD 20)
      WHEN 0 THEN 'hired' WHEN 1 THEN 'offer_accepted'
      WHEN 2 THEN 'offer_declined' WHEN 3 THEN 'offer_extended'
      WHEN 9 THEN 'withdrawn' WHEN 10 THEN 'rejected' WHEN 11 THEN 'rejected'
      WHEN 12 THEN 'rejected' WHEN 13 THEN 'rejected' WHEN 14 THEN 'rejected'
      WHEN 15 THEN 'rejected' WHEN 16 THEN 'rejected' WHEN 17 THEN 'rejected'
      WHEN 18 THEN 'rejected' WHEN 19 THEN 'rejected' ELSE 'active' END;
    SET v_stage = CASE (i MOD 20)
      WHEN 0 THEN 'hired' WHEN 1 THEN 'offer' WHEN 2 THEN 'closed'
      WHEN 3 THEN 'offer' WHEN 4 THEN 'final_interview' WHEN 5 THEN 'interview'
      WHEN 6 THEN 'phone_screen' WHEN 7 THEN 'resume_review' WHEN 8 THEN 'applied'
      WHEN 9 THEN 'closed' WHEN 10 THEN 'closed' WHEN 11 THEN 'closed'
      WHEN 12 THEN 'closed' WHEN 13 THEN 'closed' WHEN 14 THEN 'closed'
      WHEN 15 THEN 'closed' WHEN 16 THEN 'closed' WHEN 17 THEN 'closed'
      WHEN 18 THEN 'closed' ELSE 'closed' END;

    INSERT INTO applications
      (application_no, candidate_id, requisition_id, recruiter_id, applied_at,
       current_stage, status, expected_salary, available_date,
       rejection_reason, withdrawn_reason, created_at)
    VALUES
      (CONCAT('APP-2026-', LPAD(i, 4, '0')), v_candidate, v_requisition, v_recruiter,
       v_applied_at, v_stage, v_status,
       220000 + ((i MOD 12) * 25000), DATE_ADD(DATE(v_applied_at), INTERVAL (30 + (i MOD 45)) DAY),
       IF(v_status = 'rejected', ELT(((i - 1) MOD 6) + 1,
         '核心经验与岗位要求存在差距','技术深度未达到当前职级要求','岗位已优先推进其他候选人',
         '案例分析缺少关键业务指标','薪酬期望超出职位预算','综合评估后暂不匹配'), NULL),
       IF(v_status = 'withdrawn', '已接受其他公司的录用意向', NULL), v_applied_at);
    SET i = i + 1;
  END WHILE;

  INSERT INTO application_status_history
    (application_id, from_stage, to_stage, from_status, to_status, changed_by_employee_id, changed_at, change_reason)
  SELECT application_id, NULL, 'applied', NULL, 'active', recruiter_id, applied_at, '收到候选人申请'
  FROM applications;

  INSERT INTO application_status_history
    (application_id, from_stage, to_stage, from_status, to_status, changed_by_employee_id, changed_at, change_reason)
  SELECT application_id, 'applied', 'resume_review', 'active', 'active', recruiter_id,
    DATE_ADD(applied_at, INTERVAL 1 DAY), '简历进入招聘团队评审'
  FROM applications WHERE current_stage <> 'applied';

  INSERT INTO application_status_history
    (application_id, from_stage, to_stage, from_status, to_status, changed_by_employee_id, changed_at, change_reason)
  SELECT application_id, 'resume_review', 'phone_screen', 'active', 'active', recruiter_id,
    DATE_ADD(applied_at, INTERVAL 3 DAY), '简历评审通过'
  FROM applications
  WHERE current_stage IN ('phone_screen','interview','final_interview','offer','hired')
     OR (current_stage = 'closed' AND application_id MOD 3 <> 0);

  INSERT INTO application_status_history
    (application_id, from_stage, to_stage, from_status, to_status, changed_by_employee_id, changed_at, change_reason)
  SELECT application_id, 'phone_screen', 'interview', 'active', 'active', recruiter_id,
    DATE_ADD(applied_at, INTERVAL 7 DAY), '电话沟通通过，安排专业面试'
  FROM applications
  WHERE current_stage IN ('interview','final_interview','offer','hired')
     OR (current_stage = 'closed' AND application_id MOD 4 = 0);

  INSERT INTO application_status_history
    (application_id, from_stage, to_stage, from_status, to_status, changed_by_employee_id, changed_at, change_reason)
  SELECT application_id, 'interview', 'final_interview', 'active', 'active', recruiter_id,
    DATE_ADD(applied_at, INTERVAL 14 DAY), '专业面试通过'
  FROM applications WHERE current_stage IN ('final_interview','offer','hired');

  INSERT INTO application_status_history
    (application_id, from_stage, to_stage, from_status, to_status, changed_by_employee_id, changed_at, change_reason)
  SELECT application_id,
    CASE WHEN current_stage IN ('offer','hired') THEN 'final_interview' ELSE 'resume_review' END,
    current_stage, 'active', status, recruiter_id,
    DATE_ADD(applied_at, INTERVAL (18 + (application_id MOD 8)) DAY),
    CASE status WHEN 'rejected' THEN rejection_reason WHEN 'withdrawn' THEN withdrawn_reason
      WHEN 'hired' THEN '候选人已入职' WHEN 'offer_accepted' THEN '候选人接受 Offer'
      WHEN 'offer_declined' THEN '候选人拒绝 Offer' WHEN 'offer_extended' THEN 'Offer 已发送'
      ELSE '进入当前阶段' END
  FROM applications
  WHERE status <> 'active';

  SET i = 1;
  WHILE i <= 220 DO
    SELECT current_stage, status, applied_at INTO v_stage, v_status, v_applied_at
    FROM applications WHERE application_id = i;
    SET v_rounds = CASE
      WHEN v_stage IN ('offer','hired') OR v_status IN ('offer_accepted','offer_declined') THEN 3
      WHEN v_stage = 'final_interview' THEN 3
      WHEN v_stage = 'interview' THEN 2
      WHEN v_stage = 'closed' AND i MOD 4 = 0 THEN 2
      WHEN v_stage = 'closed' AND i MOD 3 <> 0 THEN 1
      ELSE 0 END;
    SET r = 1;
    WHILE r <= v_rounds DO
      SET v_mode = ELT(((i + r) MOD 3) + 1, 'video', 'onsite', 'phone');
      SET v_interview_start = DATE_ADD(v_applied_at, INTERVAL (5 + r * 5) DAY);
      SET v_interview_status = 'completed';

      IF v_status = 'active' AND r = v_rounds THEN
        SET v_interview_start = DATE_ADD('2026-07-23 09:00:00', INTERVAL ((i MOD 20) * 3 + r) HOUR);
        SET v_interview_status = IF(i MOD 2 = 0, 'confirmed', 'scheduled');
      ELSEIF v_status = 'withdrawn' AND r = v_rounds THEN
        SET v_interview_status = 'cancelled';
      ELSEIF v_status = 'rejected' AND i MOD 17 = 0 AND r = v_rounds THEN
        SET v_interview_status = 'no_show';
      END IF;

      INSERT INTO interview_sessions
        (interview_no, application_id, interview_type_id, round_number,
         scheduled_start, scheduled_end, mode, office_id, meeting_room, meeting_link,
         coordinator_id, status, candidate_confirmed_at, completed_at, cancellation_reason)
      VALUES
        (CONCAT('INT-2026-', LPAD(i, 4, '0'), '-', r), i,
         CASE r WHEN 1 THEN 1 WHEN 2 THEN IF(i MOD 3 = 0, 3, 2) ELSE IF(i MOD 2 = 0, 4, 5) END,
         r, v_interview_start, DATE_ADD(v_interview_start, INTERVAL IF(r = 1, 30, 60) MINUTE),
         v_mode, IF(v_mode = 'onsite', ((i - 1) MOD 5) + 1, NULL),
         IF(v_mode = 'onsite', CONCAT('会议室 ', CHAR(65 + (i MOD 6)), '-', 3 + (i MOD 9)), NULL),
         IF(v_mode = 'video', CONCAT('https://meet.example.test/interview/', LPAD(i, 4, '0'), '-', r), NULL),
         ((i - 1) MOD 6) + 1, v_interview_status,
         IF(v_interview_status IN ('confirmed','completed'), DATE_SUB(v_interview_start, INTERVAL 2 DAY), NULL),
         IF(v_interview_status = 'completed', DATE_ADD(v_interview_start, INTERVAL IF(r = 1, 32, 65) MINUTE), NULL),
         IF(v_interview_status = 'cancelled', '候选人已撤回申请', NULL));

      SET v_interview_id = LAST_INSERT_ID();
      SET v_interviewer_1 = 17 + ((i + r - 2) MOD 14);
      SET v_interviewer_2 = 17 + ((i + r + 5) MOD 14);
      INSERT INTO interview_panel_members
        (interview_id, employee_id, panel_role, attendance_status, assigned_at)
      VALUES
        (v_interview_id, v_interviewer_1, 'lead',
         IF(v_interview_status = 'completed', 'attended', IF(v_interview_status = 'cancelled', 'accepted', 'accepted')),
         DATE_SUB(v_interview_start, INTERVAL 4 DAY)),
        (v_interview_id, v_interviewer_2, IF(r = 3, 'bar_raiser', 'interviewer'),
         IF(v_interview_status = 'completed', 'attended', IF(i MOD 19 = 0, 'declined', 'accepted')),
         DATE_SUB(v_interview_start, INTERVAL 4 DAY));
      SET r = r + 1;
    END WHILE;
    SET i = i + 1;
  END WHILE;

  INSERT INTO interview_feedback
    (interview_id, interviewer_id, overall_rating, recommendation,
     strengths, concerns, evidence_notes, submitted_at)
  SELECT
    pm.interview_id,
    pm.employee_id,
    1 + ((pm.interview_id + pm.employee_id) MOD 5),
    ELT(1 + ((pm.interview_id + pm.employee_id) MOD 5), 'strong_no','no','mixed','yes','strong_yes'),
    ELT(1 + (pm.interview_id MOD 5),
      '问题拆解清晰，能主动识别关键约束。',
      '专业基础扎实，方案兼顾交付效率和长期维护。',
      '沟通有条理，能根据反馈及时调整方案。',
      '对业务目标敏感，能够用数据说明取舍。',
      '有较强的主人翁意识和跨团队推进经验。'),
    ELT(1 + (pm.employee_id MOD 5),
      '大规模场景下的容量评估还不够完整。',
      '回答偏重实现细节，对业务优先级讨论较少。',
      '部分关键假设未主动验证，需要提示后补充。',
      '团队管理案例的复杂度有限。',
      '暂无明显风险，建议在背调阶段核实项目职责。'),
    CONCAT('基于第 ', i.round_number, ' 轮面试记录：候选人给出了具体项目背景、个人行动和可量化结果。'),
    DATE_ADD(i.completed_at, INTERVAL (10 + (pm.employee_id MOD 80)) MINUTE)
  FROM interview_panel_members pm
  JOIN interview_sessions i ON i.interview_id = pm.interview_id
  WHERE i.status = 'completed' AND pm.attendance_status = 'attended';

  INSERT INTO feedback_competency_scores
    (feedback_id, competency_id, score, comment)
  SELECT
    f.feedback_id,
    c.competency_id,
    1 + ((f.feedback_id + c.competency_id) MOD 5),
    ELT(1 + ((f.feedback_id + c.competency_id) MOD 5),
      '缺少有效证据，表现未达到当前岗位要求。',
      '具备基础认知，但复杂场景经验仍需加强。',
      '达到岗位基本要求，能独立处理常见问题。',
      '表现稳定且有充分实例支撑，超过基本要求。',
      '展现出显著优势，可作为该能力项的团队标杆。')
  FROM interview_feedback f
  CROSS JOIN competencies c;

  INSERT INTO offers
    (offer_no, application_id, offered_title, job_level, base_salary,
     annual_bonus_target, sign_on_bonus, equity_value, proposed_start_date,
     expires_at, status, sent_at, responded_at, decline_reason, created_by_employee_id, created_at)
  SELECT
    CONCAT('OFF-2026-', LPAD(a.application_id, 4, '0')),
    a.application_id, jr.job_title, jr.job_level,
    LEAST(jr.max_salary, GREATEST(jr.min_salary, COALESCE(a.expected_salary, jr.min_salary))),
    IF(jr.job_level IN ('manager','director'), 20.00, 15.00),
    IF(a.application_id MOD 4 = 0, 30000.00, 0.00),
    IF(jr.job_level IN ('lead','manager','director'), 100000.00, 30000.00),
    DATE_ADD(DATE(a.applied_at), INTERVAL (45 + (a.application_id MOD 15)) DAY),
    DATE_ADD(a.applied_at, INTERVAL 35 DAY),
    CASE a.status WHEN 'hired' THEN 'accepted' WHEN 'offer_accepted' THEN 'accepted'
      WHEN 'offer_declined' THEN 'declined' WHEN 'offer_extended' THEN 'sent' ELSE 'approved' END,
    DATE_ADD(a.applied_at, INTERVAL 28 DAY),
    CASE WHEN a.status IN ('hired','offer_accepted','offer_declined')
      THEN DATE_ADD(a.applied_at, INTERVAL (30 + (a.application_id MOD 4)) DAY) ELSE NULL END,
    IF(a.status = 'offer_declined', ELT(1 + (a.application_id MOD 4),
      '接受了其他公司的机会','薪酬结构未达预期','工作地点与个人计划不符','希望继续留在当前公司'), NULL),
    a.recruiter_id, DATE_ADD(a.applied_at, INTERVAL 25 DAY)
  FROM applications a
  JOIN job_requisitions jr ON jr.requisition_id = a.requisition_id
  WHERE a.status IN ('hired','offer_accepted','offer_declined','offer_extended');

  INSERT INTO offer_approvals
    (offer_id, approval_step, approver_id, decision, decided_at, comment)
  SELECT offer_id, 1, 33, 'approved', DATE_ADD(created_at, INTERVAL 3 HOUR), 'HRBP 已核对职级与薪酬带宽'
  FROM offers;
  INSERT INTO offer_approvals
    (offer_id, approval_step, approver_id, decision, decided_at, comment)
  SELECT offer_id, 2, 32, 'approved', DATE_ADD(created_at, INTERVAL 8 HOUR), '预算在已批准的招聘需求范围内'
  FROM offers;
  INSERT INTO offer_approvals
    (offer_id, approval_step, approver_id, decision, decided_at, comment)
  SELECT offer_id, 3, 31,
    IF(base_salary >= 390000, 'approved', 'skipped'),
    IF(base_salary >= 390000, DATE_ADD(created_at, INTERVAL 1 DAY), NULL),
    IF(base_salary >= 390000, '关键岗位薪酬审批通过', '未达到 VP 审批阈值')
  FROM offers;

  UPDATE job_requisitions jr
  SET filled_count = LEAST(headcount, (
      SELECT COUNT(*) FROM applications a
      WHERE a.requisition_id = jr.requisition_id AND a.status = 'hired'
    ));
END$$

CALL seed_interview_db()$$
DROP PROCEDURE seed_interview_db$$

DELIMITER ;

ANALYZE TABLE
  candidates, applications, interview_sessions, interview_feedback,
  feedback_competency_scores, offers;
