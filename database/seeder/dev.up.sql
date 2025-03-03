DO
$$
    DECLARE
        -- User IDs declarations
        student1_id UUID;
        student2_id UUID;
        student3_id UUID;
        mentor1_id  UUID;
        mentor2_id  UUID;
        admin_id    UUID;

        -- Category IDs declarations (UUIDv7 format)
        webdev_id        UUID;
        datascience_id   UUID;
        mobile_id        UUID;
        cloud_id         UUID;
        ai_id            UUID;
        devops_id        UUID;
        cybersecurity_id UUID;

        -- Course IDs declarations (UUIDv7 format)
        course1_id       UUID;
        course2_id       UUID;
        course3_id       UUID;
        course4_id       UUID;
        course5_id       UUID;
        course6_id       UUID;
        course7_id       UUID;
        course8_id       UUID;
        course9_id       UUID;
        course10_id      UUID;
        course11_id      UUID;
        course12_id      UUID;
        course13_id      UUID;
        course14_id      UUID;
        course15_id      UUID;
        course16_id      UUID;
        course17_id      UUID;
        course18_id      UUID;
        course19_id      UUID;
        course20_id      UUID;
        course21_id      UUID;
        course22_id      UUID;
        course23_id      UUID;
        course24_id      UUID;
        course25_id      UUID;
    BEGIN
        -- Initialize User UUIDs
        student1_id := '0194bb8e-1e7c-4082-806a-1c9483a59a1b';
        student2_id := '0194bbc5-0cfc-408e-802d-a854785c39d9';
        student3_id := '0194bbfb-fb7c-40d9-8029-7ca996ab5da6';
        mentor1_id := '0194bc32-e9fc-405c-801a-08ff3b0cf28b';
        mentor2_id := '0194bc69-d87c-40d5-809a-17513f3d2b98';
        admin_id := '0194bca0-c6fc-40ec-8051-68235bef9817';

        -- Initialize Category UUIDs (UUIDv7 format)
        webdev_id := '01890e5c-b334-7c38-8c8c-b05c57968489';
        datascience_id := '01890e5c-bed5-7ebb-8dcf-ada8c40c2253';
        mobile_id := '01890e5c-c5b6-7da1-b35c-70d83a37d078';
        cloud_id := '01890e5c-cdbe-718c-9623-ec4a1bdfaa27';
        ai_id := '01890e5c-d583-752e-ab1e-24d75a5c4dfc';
        devops_id := '01890e5c-dc95-7de2-9e36-00a6a28bf8de';
        cybersecurity_id := '01890e5c-e31e-7b65-b6fb-08a3b7a00d56';

        -- Initialize Course UUIDs (UUIDv7 format)
        course1_id := '01890e5d-05a6-7b0b-a7b3-3f8ed37b5261';
        course2_id := '01890e5d-0ec9-7e29-ae5f-4a7c66f5d456';
        course3_id := '01890e5d-1657-71ae-bcba-54a2a55b1bb8';
        course4_id := '01890e5d-1e20-7e5c-b0b2-6dbc90efac62';
        course5_id := '01890e5d-25de-7d73-bf1c-8c1a13f17839';
        course6_id := '01890e5d-2dd5-7bde-994a-76e4b9e18a08';
        course7_id := '01890e5d-3530-7487-8b77-fc5cc0773b3f';
        course8_id := '01890e5d-3c4a-72fe-b78b-e1ebf0cbed87';
        course9_id := '01890e5d-437c-72ca-8ddc-4aa6b26ddc11';
        course10_id := '01890e5d-4a4c-7452-b55a-7b8fa75d84e5';
        course11_id := '01890e5d-516c-7dc0-9d9c-1d6ec6e36fc4';
        course12_id := '01890e5d-58ca-787d-9ed5-31fbe2bc9429';
        course13_id := '01890e5d-6086-7c90-a81b-9dccca9a6c97';
        course14_id := '01890e5d-6769-7b84-baba-4b3aad93a8b0';
        course15_id := '01890e5d-6eda-71bf-8a5e-a74a7f8c8d2c';
        course16_id := '01890e5d-76d1-75af-9c61-0adf9a0d8db1';
        course17_id := '01890e5d-7e5d-7bfb-9e26-2eb534b456f3';
        course18_id := '01890e5d-8540-7a65-9b7d-5d95efed8bc9';
        course19_id := '01890e5d-8c74-7b15-9ed8-15b5e14ab0c7';
        course20_id := '01890e5d-946d-7a56-a9b4-69d7c6fefb92';
        course21_id := '01890e5d-9b95-72e2-a064-a37c64ac6d3a';
        course22_id := '01890e5d-a316-7e25-bad9-0ec1a4eebf79';
        course23_id := '01890e5d-aa41-7ec4-9278-a1eaf2fa3b40';
        course24_id := '01890e5d-b181-7f2d-aeda-c4e69b0b3d75';
        course25_id := '01890e5d-b8cf-7d50-b3a1-1acf37eff9af';

        -- Users seeder
        INSERT INTO users (id, name, email, password_hash, role, has_avatar, created_at, updated_at)
        VALUES
            -- Students
            (student1_id, 'Student 1', 'student1@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'student',
             FALSE,
             NOW() - INTERVAL '30 days',
             NOW() - INTERVAL '30 days'),

            (student2_id, 'Student 2', 'student2@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'student',
             FALSE,
             NOW() - INTERVAL '30 days' + INTERVAL '1 hour',
             NOW() - INTERVAL '30 days' + INTERVAL '1 hour'),

            (student3_id, 'Student 3', 'student3@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'student',
             FALSE,
             NOW() - INTERVAL '30 days' + INTERVAL '2 hours',
             NOW() - INTERVAL '30 days' + INTERVAL '2 hours'),

            -- Mentors
            (mentor1_id, 'Mentor 1', 'mentor1@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'mentor',
             FALSE,
             NOW() - INTERVAL '30 days' + INTERVAL '3 hours',
             NOW() - INTERVAL '30 days' + INTERVAL '3 hours'),

            (mentor2_id, 'Mentor 2', 'mentor2@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'mentor',
             FALSE,
             NOW() - INTERVAL '30 days' + INTERVAL '4 hours',
             NOW() - INTERVAL '30 days' + INTERVAL '4 hours'),

            -- Admin
            (admin_id, 'Admin 1', 'admin1@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'admin',
             FALSE,
             NOW() - INTERVAL '30 days' + INTERVAL '5 hours',
             NOW() - INTERVAL '30 days' + INTERVAL '5 hours');

        -- Students specific data
        INSERT INTO students (user_id, instance, major)
        VALUES (student1_id, 'University of Technology', 'Computer Science'),
               (student2_id, 'State University', 'Data Science'),
               (student3_id, 'National Institute', 'Software Engineering');

        -- Mentors specific data
        INSERT INTO mentors (user_id, specialization, experience, rating, rating_count, rating_total, price, balance)
        VALUES (mentor1_id, 'Web Development',
                'Over 8 years of experience in full-stack development with expertise in React, Node.js, and database design. Led development teams for enterprise applications and mentored junior developers.',
                4.8, 125, 4.8 * 125, 750000, 3500000),

               (mentor2_id, 'Machine Learning',
                'Data scientist with 6 years of experience in implementing machine learning models for production. Background in computer vision and natural language processing. Published research in top AI conferences.',
                4.7, 98, 4.7 * 98, 850000, 2750000);

        -- Categories seeder
        INSERT INTO categories (id, name)
        VALUES (webdev_id, 'Web Development'),
               (datascience_id, 'Data Science'),
               (mobile_id, 'Mobile Development'),
               (cloud_id, 'Cloud Computing'),
               (ai_id, 'Artificial Intelligence'),
               (devops_id, 'DevOps'),
               (cybersecurity_id, 'Cybersecurity');

        -- Courses seeder
        INSERT INTO courses (id, category_id, title, description, teacher_name, rating, rating_count, total_rating,
                             enrollment_count, content_count, created_at, updated_at)
        VALUES
            -- Web Development courses
            (course1_id,
             webdev_id,
             'Modern JavaScript Fundamentals',
             'Master the core concepts of JavaScript including ES6+ features, asynchronous programming, and functional programming techniques. This comprehensive course covers everything from basic syntax to advanced topics like closures, promises, and modern tooling.',
             'Dr. Sarah Johnson',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '60 days',
             NOW() - INTERVAL '15 days'),

            (course2_id,
             webdev_id,
             'Full-Stack React & Node.js',
             'Build production-ready applications with React and Node.js. Learn state management with Redux, server-side rendering, RESTful API design, authentication, testing, and deployment strategies for modern web applications.',
             'Michael Chen',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '45 days',
             NOW() - INTERVAL '10 days'),

            (course3_id,
             webdev_id,
             'CSS Mastery: Advanced Styling',
             'Take your CSS skills to the next level with advanced selectors, animations, grid layouts, and CSS variables. Learn responsive design techniques, CSS architecture patterns, and optimization strategies for modern websites.',
             'Sophia Martinez',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '85 days',
             NOW() - INTERVAL '12 days'),

            (course4_id,
             webdev_id,
             'Vue.js for Frontend Development',
             'Learn Vue.js from scratch and build dynamic user interfaces. This course covers Vue components, directives, Vuex for state management, Vue Router for navigation, and integration with backend APIs.',
             'David Wilson',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '70 days',
             NOW() - INTERVAL '8 days'),

            -- Data Science courses
            (course5_id,
             datascience_id,
             'Python for Data Analysis',
             'Learn how to use Python libraries like Pandas, NumPy, and Matplotlib to analyze, visualize, and interpret complex datasets. This course covers data cleaning, transformation, exploratory analysis, and creating insightful visualizations.',
             'Dr. Emily Rodriguez',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '90 days',
             NOW() - INTERVAL '25 days'),

            (course6_id,
             datascience_id,
             'Statistical Methods for Data Science',
             'Develop a strong foundation in statistics essential for data science. Topics include probability distributions, hypothesis testing, regression analysis, experimental design, and Bayesian methods with practical applications using Python and R.',
             'Prof. James Wilson',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '75 days',
             NOW() - INTERVAL '20 days'),

            (course7_id,
             datascience_id,
             'Big Data Processing with Spark',
             'Master distributed data processing using Apache Spark. Learn to handle large-scale datasets, perform batch and stream processing, and implement machine learning pipelines in a distributed environment.',
             'Dr. Aisha Patel',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '65 days',
             NOW() - INTERVAL '18 days'),

            (course8_id,
             datascience_id,
             'Data Visualization with D3.js',
             'Create interactive and dynamic data visualizations for the web using D3.js. Learn how to transform data into compelling visual stories through charts, graphs, and interactive dashboards.',
             'Marcus Johnson',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '55 days',
             NOW() - INTERVAL '14 days'),

            -- Mobile Development courses
            (course9_id,
             mobile_id,
             'iOS App Development with Swift',
             'Create beautiful and functional iOS applications using Swift and UIKit. Learn app architecture, interface design principles, Core Data, networking, and publishing to the App Store. Build several complete applications throughout the course.',
             'Alex Turner',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '120 days',
             NOW() - INTERVAL '30 days'),

            (course10_id,
             mobile_id,
             'Android Development with Kotlin',
             'Build modern Android applications using Kotlin and Jetpack components. Master Android architecture, UI design with Material Design, data persistence, background processing, and Google Play deployment.',
             'Priya Sharma',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '110 days',
             NOW() - INTERVAL '28 days'),

            (course11_id,
             mobile_id,
             'Cross-Platform Apps with Flutter',
             'Develop cross-platform mobile applications with a single codebase using Flutter and Dart. Learn to create beautiful UI components, manage state, integrate APIs, and deploy to both iOS and Android platforms.',
             'Ryan Park',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '100 days',
             NOW() - INTERVAL '26 days'),

            -- Cloud Computing courses
            (course12_id,
             cloud_id,
             'AWS Solutions Architect',
             'Prepare for the AWS Solutions Architect certification while learning to design distributed systems on AWS. Cover EC2, S3, RDS, Lambda, API Gateway, CloudFormation, and best practices for security, cost optimization, and high availability.',
             'Jennifer Parker',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '100 days',
             NOW() - INTERVAL '5 days'),

            (course13_id,
             cloud_id,
             'Google Cloud Platform Essentials',
             'Learn the fundamentals of Google Cloud Platform services and architecture. Explore Compute Engine, App Engine, Cloud Storage, BigQuery, Cloud Functions, and deployment automation with hands-on projects.',
             'Thomas Lee',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '95 days',
             NOW() - INTERVAL '22 days'),

            (course14_id,
             cloud_id,
             'Microsoft Azure Administration',
             'Master Azure infrastructure management including virtual machines, storage, networking, security, and monitoring. Prepare for the Azure Administrator certification with practical exercises and real-world scenarios.',
             'Samantha Wright',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '85 days',
             NOW() - INTERVAL '16 days'),

            (course15_id,
             cloud_id,
             'Serverless Architecture Patterns',
             'Design and implement scalable serverless applications across major cloud providers. Learn event-driven architecture, function-as-a-service, managed services integration, and deployment automation for serverless systems.',
             'Daniel Garcia',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '75 days',
             NOW() - INTERVAL '10 days'),

            -- AI courses
            (course16_id,
             ai_id,
             'Machine Learning Fundamentals',
             'Build a strong foundation in machine learning algorithms and techniques. This course covers supervised and unsupervised learning, model evaluation, feature engineering, and practical implementation using Python''s scikit-learn library.',
             'Dr. Robert Kim',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '85 days',
             NOW() - INTERVAL '12 days'),

            (course17_id,
             ai_id,
             'Deep Learning with TensorFlow',
             'Dive into neural networks and deep learning using TensorFlow and Keras. Learn to build, train and deploy convolutional neural networks, recurrent neural networks, and transformers for computer vision, NLP, and other applications.',
             'Dr. Lisa Zhang',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '70 days',
             NOW() - INTERVAL '8 days'),

            (course18_id,
             ai_id,
             'Natural Language Processing',
             'Master techniques for processing and analyzing text data using Python. Learn text preprocessing, feature extraction, sentiment analysis, topic modeling, and building language models with transformers like BERT and GPT.',
             'Prof. Alan Foster',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '60 days',
             NOW() - INTERVAL '15 days'),

            (course19_id,
             ai_id,
             'Computer Vision Applications',
             'Learn to build applications that can interpret and process visual information from images and videos. Cover image classification, object detection, facial recognition, and generative models using PyTorch and OpenCV.',
             'Dr. Nina Zhou',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '50 days',
             NOW() - INTERVAL '7 days'),

            -- DevOps courses
            (course20_id,
             devops_id,
             'Docker and Kubernetes for DevOps',
             'Master containerization with Docker and orchestration with Kubernetes. Learn to build, deploy, and scale microservices architectures using container technologies with CI/CD pipeline integration.',
             'Richard Stevens',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '90 days',
             NOW() - INTERVAL '14 days'),

            (course21_id,
             devops_id,
             'CI/CD with Jenkins and GitLab',
             'Implement continuous integration and deployment pipelines using Jenkins and GitLab. Learn to automate building, testing, and deploying applications with version control integration and quality assurance.',
             'Michelle Carter',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '80 days',
             NOW() - INTERVAL '18 days'),

            (course22_id,
             devops_id,
             'Infrastructure as Code with Terraform',
             'Manage infrastructure declaratively using Terraform across multiple cloud providers. Learn to define, provision, and version control infrastructure resources with automated deployment workflows.',
             'Jason Miller',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '70 days',
             NOW() - INTERVAL '10 days'),

            -- Cybersecurity courses
            (course23_id,
             cybersecurity_id,
             'Ethical Hacking Fundamentals',
             'Learn the methodology and tools used by ethical hackers to identify vulnerabilities in systems and networks. Practice penetration testing techniques, vulnerability assessment, and secure system hardening.',
             'Carlos Mendoza',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '85 days',
             NOW() - INTERVAL '20 days'),

            (course24_id,
             cybersecurity_id,
             'Network Security Architecture',
             'Design and implement secure network infrastructures to protect organizational assets. Cover firewalls, intrusion detection systems, VPNs, access control, and security protocols with practical defense strategies.',
             'Dr. Olivia Thompson',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '75 days',
             NOW() - INTERVAL '15 days'),

            (course25_id,
             cybersecurity_id,
             'Secure Coding Practices',
             'Learn to write secure code that prevents common vulnerabilities like injection attacks, cross-site scripting, and authentication flaws. Implement security controls and code review techniques across different programming languages.',
             'Benjamin Okon',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '65 days',
             NOW() - INTERVAL '8 days');
    END
$$;
