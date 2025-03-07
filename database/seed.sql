DELETE FROM recipes;
INSERT INTO recipes (name, user_id, description, cooking_time, instructions) VALUES 
('Garlic Parmesan Chicken', 
 1,
 'A savory dish of chicken breast seasoned with garlic and parmesan cheese.',
 '45 minutes',
 'Preheat oven to 400°F (200°C). Season chicken breasts with salt and garlic. Heat olive oil in a skillet and sear chicken until golden. Transfer to a baking dish, top with Parmesan cheese. Bake for 25-30 minutes until chicken is cooked through.'
),
('Spaghetti with Tomato Basil Sauce', 
 1,
 'Classic Italian spaghetti served with a flavorful tomato and basil sauce.',
 '30 minutes',
 'Cook spaghetti according to package instructions, then drain. Heat olive oil in a saucepan, sauté garlic until fragrant. Add chopped tomatoes and basil, simmer for 15 minutes. Mix cooked spaghetti with the sauce. Serve hot and garnish with additional basil if desired.'
),
('Chocolate Chip Pancakes', 
 1,
 'Fluffy pancakes loaded with gooey chocolate chips.',
 '20 minutes',
 'In a large bowl, mix flour, sugar, salt, and baking powder. Whisk in eggs, milk, and melted butter until smooth. Heat a non-stick pan over medium heat. Pour batter onto the pan and sprinkle with chocolate chips. Cook until bubbles form, then flip and cook until golden brown. Serve warm with your favorite toppings.'
),
('Zucchini Noodles with Pesto and Grilled Chicken', 
 1, 
 '20 minutes', 
 'A tasty low-carb pasta dish',
 'Heat a grill pan and cook the chicken breasts until fully cooked. Set aside to rest. In a food processor, blend basil, garlic, pine nuts, parmesan cheese, and olive oil until smooth to make the pesto. Use a spiralizer to create zucchini noodles. Heat a skillet with olive oil, toss in the zucchini noodles for 2-3 minutes until slightly tender. Mix the pesto into the noodles and plate them. Slice the grilled chicken and place it on top of the noodles. Serve warm and enjoy.');


DELETE FROM ingredients;
-- Garlic Parmesan Chicken
INSERT INTO ingredients (name, user_id, recipe_id, quantity, unit) VALUES
('chicken breast', 1, 1, 2, 'pcs'),    
('garlic', 1, 1, 3, 'cloves'),  
('olive oil', 1, 1, 2, 'tbsp'),     
('parmesan cheese', 1, 1, 0.5, 'cup');

-- Spaghetti with Tomato Basil Sauce
INSERT INTO ingredients (name, user_id, recipe_id, quantity, unit) VALUES
('spaghetti', 1, 2, 200, 'grams'),
('tomatoes', 1, 2, 3, 'pcs'), 
('basil', 1, 2, 5, 'leaves'),  
('olive oil', 1, 2, 1, 'tbsp'),  
('garlic', 1, 2, 1, 'clove'); 

-- Chocolate Chip Pancakes
INSERT INTO ingredients (name, user_id, recipe_id, quantity, unit) VALUES
('flour', 1, 3, 1.5, 'cups'),   
('sugar', 1, 3, 0.25, 'cup'),
('eggs', 1, 3, 2, 'pcs'),      
('milk', 1, 3, 1, 'cup'),  
('butter', 1, 3, 2, 'tbsp'),      
('salt', 1, 3, 0.5, 'tsp');    

INSERT INTO ingredients (name, user_id, recipe_id, quantity, unit) VALUES
('chicken breast', 1, 4, 2, 'piece'),         
('zucchini', 1, 4, 3, 'piece'),        
('basil', 1, 4, 1, 'cup'),             
('garlic', 1, 4, 2, 'clove'),     
('pine nuts', 1, 4, 0.25, 'cup'),  
('parmesan cheese', 1, 4, 0.5, 'cup'),      
('olive oil', 1, 4, 3, 'tbsp'),      
('salt', 1, 4, 0.5, 'tsp'),
('pepper', 1, 4, 0.25, 'tsp');  

INSERT INTO recipe_tags (recipe_id, name) VALUES 
(2, 'vegetarian'),
(3, 'vegetarian'),
(4, 'low-carb');


SELECT * FROM recipes;
SELECT * FROM ingredients;
