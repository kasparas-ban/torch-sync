INSERT INTO items (user_id, public_item_id, title, type_, priority) VALUES (1, "1ax1usfu2uku", "Learn Spanish", "DREAM", "HIGH");
INSERT INTO items (user_id, public_item_id, title, type_) VALUES (1, "2ax1usfu2uku", "Get fit", "DREAM");
INSERT INTO items (user_id, public_item_id, title, type_) VALUES (1, "3ax1usfu2uku", "Get good at math", "DREAM");

INSERT INTO items (user_id, public_item_id, title, type_, target_date, priority) VALUES (1, "4bax1usfu2uk", "Make a todo/timer app", "GOAL", "2023-12-01", "HIGH");
INSERT INTO items (user_id, public_item_id, title, type_, priority) VALUES (1, "5bax1usfu2uk", "Learn chess", "GOAL", "LOW");
INSERT INTO items (user_id, public_item_id, title, type_, parent_id) VALUES (1, "6bax1usfu2uk", "Learn Spanish vocabulary", "GOAL", "1ax1usfu2uku");
INSERT INTO items (user_id, public_item_id, title, type_, parent_id) VALUES (1, "7bax1usfu2uk", "Learn Spanish grammar", "GOAL", "1ax1usfu2uku");
INSERT INTO items (user_id, public_item_id, title, type_, parent_id) VALUES (1, "8bax1usfu2uk", "Spanish language comprehension", "GOAL", "1ax1usfu2uku");
INSERT INTO items (user_id, public_item_id, title, type_, parent_id) VALUES (1, "9bax1usfu2uk", "Spanish writing", "GOAL", "1ax1usfu2uku");
INSERT INTO items (user_id, public_item_id, title, type_, parent_id) VALUES (1, "10ax1usfu2uk", "Build muscle", "GOAL", "2ax1usfu2uku");
INSERT INTO items (user_id, public_item_id, title, type_, parent_id, target_date) VALUES (1, "11ax1usfu2uk", "Learn Linear Algebra", "GOAL", "3ax1usfu2uku", "2023-12-01");
INSERT INTO items (user_id, public_item_id, title, type_, parent_id) VALUES (1, "12ax1usfu2uk", "Learn Calculus", "GOAL", "3ax1usfu2uku");
INSERT INTO items (user_id, public_item_id, title, type_) VALUES (1, "13ax1usfu2uk", "Read "Demons" by Dostoevsky", "GOAL");
INSERT INTO items (user_id, public_item_id, title, type_) VALUES (1, "14ax1usfu2uk", "Read "The Shape of Space"", "GOAL");

INSERT INTO items (user_id, public_item_id, title, type_, duration, time_spent, target_date, priority, parent_id) VALUES (1, "15ax1usfu2uk", "Make a Figma design sketch", "TASK", 100800, 90000, "2023-10-30", "MEDIUM", "4bax1usfu2uk");
INSERT INTO items (user_id, public_item_id, title, type_, duration, time_spent, target_date, priority, parent_id) VALUES (1, "16ax1usfu2uk", "Code MVP frontend", "TASK", 144000, 85000, "2023-10-30", "HIGH", "4bax1usfu2uk");
INSERT INTO items (user_id, public_item_id, title, type_, duration, time_spent, target_date, priority, parent_id) VALUES (1, "17ax1usfu2uk", "Make MVP backend", "TASK", 108000, 40000, "2023-10-30", "HIGH", "4bax1usfu2uk");
INSERT INTO items (user_id, public_item_id, title, type_, duration, time_spent, parent_id) VALUES (1, "18ax1usfu2uk", "Learn common Spanish greeting phrases", "TASK", 36000, 8000, "6bax1usfu2uk");
INSERT INTO items (user_id, public_item_id, title, type_, duration, time_spent, parent_id) VALUES (1, "19ax1usfu2uk", "Memorize a list of essential words", "TASK", 36000, 1000, "6bax1usfu2uk");
INSERT INTO items (user_id, public_item_id, title, type_, duration, time_spent, parent_id) VALUES (1, "20ax1usfu2uk", "Learn Spanish pronunciation", "TASK", 36000, 30000, "6bax1usfu2uk");
INSERT INTO items (user_id, public_item_id, title, type_, parent_id, rec_times, rec_period, rec_progress) VALUES (1, "21ax1usfu2uk", "Do weight lifting", "TASK", "10ax1usfu2uk", 3, "WEEK", 2);