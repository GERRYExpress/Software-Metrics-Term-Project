INSERT INTO accounts (username, password, created_date, is_seller) VALUES ('gerry', '123', now(), true);


INSERT INTO categories (name) VALUES ('Dark Fantasy'), ('Atmospheric'), ('RPG'), ('Adventure'),
('Strategy'),
('Simulation'),
('First-Person Shooter'),
('Horror'),
('Open World'),
('Sports'),
('Racing'),
('Indie'),
('Platformer'),
('Massively Multiplayer'),
('Fighting'),
('Puzzle');

INSERT INTO products (name,description,price, poster_url)
VALUES
('DARK SOULS™ III', 'Dark Souls continues to push the boundaries with the latest, ambitious chapter in the critically-acclaimed and genre-defining series. Prepare yourself and Embrace The Darkness!', 59.99, 'https://cdn2.steamgriddb.com/grid/55e39b04fb5ce386d690685f44029957-fakepng.png'),
('DARK SOULS™: REMASTERED', 'Then, there was fire. Re-experience the critically acclaimed, genre-defining game that started it all. Beautifully remastered, return to Lordran in stunning high-definition detail running at 60fps.',  39.99, 'https://cdn2.steamgriddb.com/grid/d68159bcc8608831094ffeac4bb8b5a1-fakepng.png'),
('DARK SOULS™ II: Scholar of the First Sin', 'DARK SOULS™ II: Scholar of the First Sin brings the franchise''s renowned obscurity & gripping gameplay to a new level. Join the dark journey and experience overwhelming enemy encounters, diabolical hazards, and unrelenting challenge.', 39.99, 'https://cdn2.steamgriddb.com/grid/3c55b687e66336ed79466d9d32f14058.png'),
('ELDEN RING', 'THE CRITICALLY ACCLAIMED FANTASY ACTION RPG. Rise, Tarnished, and be guided by grace to brandish the power of the Elden Ring and become an Elden Lord in the Lands Between.', 59.99, 'https://cdn2.steamgriddb.com/grid/2ef2c21c3f541cb0126c325ee5d33864-fakepng.png'),
('ELDEN RING NIGHTREIGN', 'ELDEN RING NIGHTREIGN is a standalone adventure within the ELDEN RING universe, crafted to offer players a new gaming experience by reimagining the game''s core design.', 39.99, 'https://cdn2.steamgriddb.com/grid/8a1cef76c3ea963a331d1cebea6ef7d2.jpg'),
('ELDEN RING Shadow of the Erdtree', 'The ELDEN RING Shadow of the Erdtree expansion features an all-new story set in the Land of Shadow imbued with mystery, perilous dungeons, and new enemies, weapons and equipment.', 39.99, 'https://cdn2.steamgriddb.com/grid/73d5e7c59ebbbaf1fef3850dee8c6a1f.png'),
('Bloodborne', 'Hunt your nightmares A lone traveler. A cursed town. A deadly mystery that swallows everything it touches. Face your fears as you enter the decaying city of Yharnam, a forsaken place ravaged by a terrible, all-consuming illness. Scour its darkest shadows, fight for your life with blades and guns and discover secrets that will make your blood run cold – but just might save your skin...', 19.99, 'https://cdn2.steamgriddb.com/grid/99ce11be28659ce5cb0a3ec79a5f06df-fakepng.png'),
('Sekiro™: Shadows Die Twice - GOTY Edition', 'Game of the Year - The Game Awards 2019 Best Action Game of 2019 - IGN Carve your own clever path to vengeance in the award winning adventure from developer FromSoftware, creators of Bloodborne and the Dark Souls series. Take Revenge. Restore Your Honor. Kill Ingeniously.', 59.99, 'https://cdn2.steamgriddb.com/grid/500f4c19df482dc68975711067cb48af-fakepng.png'),
('Demon''s Souls',          'The legend that started it all, rebuilt from the ground up. Journey to the cursed kingdom of Boletaria, a once-proud land consumed by a dark scourge. Face fearsome demons and insurmountable odds as you unravel the mystery behind the fog.',59.99, 'https://cdn2.steamgriddb.com/grid/d295f66b5c911d81c5f4412a821cbedf.png'),
('Lies of P',               'Lies of P is a dark and grim soulslike game based on the familiar story of Pinocchio. Set in the Belle Époque-inspired city of Krat, you play as P, a puppet searching for Mr. Geppetto amidst mechanical puppet madness.',                                                                              49.99, 'https://cdn2.steamgriddb.com/grid/dc6e8312db8058502e345fdd2eb70d9f.png'),
('Lords of the Fallen',     'A vast world awaits in the all-new Lords of the Fallen. As one of the fabled Dark Crusaders, embark on an epic quest to overthrow Adyr, the demon God, across both the worlds of the living and the dead.',                                                                                             59.99, 'https://cdn2.steamgriddb.com/grid/19cca09485ba2f1c1fb4ec41b5e48a0f.jpg'),
('Hollow Knight',           'Forge your own path in Hollow Knight, a challenging 2D action-adventure. Journey through vast underground kingdoms, fight enemies and bosses, and uncover an ancient mystery at the heart of a long-dead civilization.',                                                                               14.99, 'https://cdn2.steamgriddb.com/grid/2cffdc4195ce6adf0a57062e4318662e-fakepng.png'),
('Hollow Knight: Silksong', 'Journey to an all-new kingdom, built upward into the clouds. Play as Hornet, princess-protector of Hallownest, and ascend the citadel of Pharloom through a teeming world of insects and wonders.',                                                                                                 29.99, 'https://cdn2.steamgriddb.com/grid/7282feea57f0e6cd89092df704f8e039-fakepng.png'),
('Blasphemous',             'Blasphemous is a brutal action-platformer set in the nightmare world of Cvstodia. Play as The Penitent One and fight your way through this land of twisted religion while trying to discover the origin of The Miracle.',                                                                              24.99, 'https://cdn2.steamgriddb.com/grid/25c2c7b955a7e2580daa62534486ec09.jpg'),
('Blasphemous 2',           'The Penitent One awakens once more. A new miracle has been foretold — join the battle against the Three Penitences and face hordes of sinful creatures in a beautifully hand-drawn dark fantasy world.',                                                                                              29.99, 'https://cdn2.steamgriddb.com/grid/f448516500da6c00b831b48426579f7e.jpg'),
('Code Vein',               'In the not too distant future, a mysterious disaster has ravaged the world as we know it. You are a Revenant, an undead hunter known as a Vein. Lost in a crumbling world, band together with companions to find a way out.',                                                                        39.99, 'https://cdn2.steamgriddb.com/grid/34353396d7f35dc07c1c7fee153932f4.png'),
('Mortal Shell',            'Mortal Shell is a deep action-RPG that tests your sanity and resilience in a shattered world. As a foundling, you discover the power to inhabit and possess the decaying remains of defeated warriors.',                                                                                              29.99, 'https://cdn2.steamgriddb.com/grid/3b6d20b48a1f1fc814e7aaf9e5bff5b1.jpg'),
('Nioh 2 – Complete Edition','Master the lethal arts of the samurai in this brutal action RPG. A prequel to the original Nioh, Nioh 2 features enhanced masocore elements and a dark, supernatural backdrop of a war-ravaged Japan.',                                                                                               49.99, 'https://cdn2.steamgriddb.com/grid/50502027c1db52caa91b31b7a7349c1a.png'),
('Wo Long: Fallen Dynasty',  'A dark fantasy Three Kingdoms period tale of a nameless militia soldier. Experience intense battles and achieve heroic feats as you rise up through the midst of the Yellow Turban Rebellion.',                                                                                                      59.99, 'https://cdn2.steamgriddb.com/grid/77f2de998d291ded2573171c666b5565.jpg'),
('Ghost of Tsushima',       'In the late 13th century, the Mongol empire has laid waste to the island of Tsushima. As one of the last surviving samurai, you rise from the wreckage and take on the role of the Ghost, an unconventional warrior.',                                                                               59.99, 'https://cdn2.steamgriddb.com/grid/0b804a2dd7d91a80df74bc3f7680ec77.png'),
('Dragon''s Dogma 2',       'A single-player focused action-RPG that pushes the boundaries of the genre. Experience a rich world full of mythical beasts and daunting foes in a land that awaits exploration.',                                                                                                                   59.99, 'https://cdn2.steamgriddb.com/grid/8cab2185b82b866fbcda56d946276f98.jpg'),
('The Witcher 3: Wild Hunt', 'You are Geralt of Rivia, a professional monster-hunter known as a Witcher. You have a contract to find a child of prophecy — one who has the power to stop the White Frost and save the world.',                                                                                                     39.99, 'https://cdn2.steamgriddb.com/grid/81d7c3d7086dd2365d06792a2f732ba6.png'),
('Cyberpunk 2077',          'Cyberpunk 2077 is an open-world action adventure set in Night City, a megalopolis obsessed with power, glamour and body modification. You play as V, a mercenary outlaw going after a one-of-a-kind implant.',                                                                                        59.99, 'https://cdn2.steamgriddb.com/grid/6004cb7860684c1cc497d133e1dae5ca-fakepng.png'),
('Baldur''s Gate 3',        'Gather your party and return to the Forgotten Realms in a tale of fellowship and betrayal, sacrifice and survival, and the lure of absolute power. An epic RPG crafted with an unprecedented level of freedom.',                                                                                      59.99, 'https://cdn2.steamgriddb.com/grid/2ddd2ff467042c2d0d238ddfbc77c0fb-fakepng.png'),
('Divinity: Original Sin 2', 'The Divine is dead. The Void approaches. The battle for Divinity has begun. Gather your companions and pursue your path to divinity in an epic story of conflict and power in Rivellon.',                                                                                                           44.99, 'https://cdn2.steamgriddb.com/grid/0f530b1c3f33c6409f7e4f10ca8e3e96.jpg'),
('Hades','Defy the god of the dead as you hack and slash your way out of the Underworld in this rogue-like dungeon crawler from the creators of Bastion, Transistor, and Pyre.',24.99, 'https://cdn2.steamgriddb.com/grid/01f34dfe5d0c4e5aa25d00beea480eab.png'),
('Hades II',                'The sequel to the award-winning rogue-like dungeon crawler. Play as the immortal princess of the Underworld and take the fight to the Titan of Time in a darker, deeper mythological adventure.',                                                                                                    29.99, 'https://cdn2.steamgriddb.com/grid/69f493afec6926ca2b60938c646bb765-fakepng.png'),
('Dead Cells',              'You are dead. You find yourself in a dungeon, reanimated by a mysterious slime. Hack, slash, and parry your way through procedurally-generated dungeons in this rogue-lite, Metroidvania-inspired action-platformer.',                                                                                24.99, 'https://cdn2.steamgriddb.com/grid/0816d3a34b9f13e7e964d4d1d83aa505.png'),
('Salt and Sanctuary',      'A 2D action RPG set in a grim and foreboding world. On a remote island you must battle grotesque bosses, collect loot, and level up your character using skill trees and an intricate system of weapon upgrades.',                                                                                    17.99, 'https://cdn2.steamgriddb.com/grid/b4858155e91937d6a75b19ce0e2ffe1f.png'),
('Remnant II',              'Remnant II is the sequel to the best-kept secret in games. Lost in a world swallowed by darkness, you and up to two friends will work together to survive hordes of deadly enemies and god-like bosses.',                                                                                             49.99, 'https://cdn2.steamgriddb.com/grid/848708d1266ecca55c4272e26ce36f50.png'),
('Star Wars Jedi: Survivor', 'Cal Kestis continues his journey in the sequel to Star Wars Jedi: Fallen Order. Travel to mysterious new planets and confront your past in this action-adventure set a galaxy far, far away.',                                                                                                       59.99, 'https://cdn2.steamgriddb.com/grid/08c38940fb115f00f5d0dac46e4ac639.png'),
('Final Fantasy XVI',       'Darkness covers the realm of Valisthea. Magic is fading, and powerful nations crumble as they wage war for the last remaining Mothercrystals. One young man must rise as a warrior and seek the truth behind the calamity.',                                                                          59.99, 'https://cdn2.steamgriddb.com/grid/528389d3fa104d48cf5ca7c97cc1b67f-fakepng.png'),
('God of War',              'His vengeance against the Gods of Olympus behind him, Kratos now lives in the Norse realm of Midgard. But past haunts return as a new journey begins — one that will change him and his son forever.',                                                                                                49.99, 'https://cdn2.steamgriddb.com/grid/5855660034a74cfe0e5fc8d57d17f4ac.png'),
('God of War Ragnarök',     'Fimbulwinter is upon them. Kratos and Atreus must journey to each of the Nine Realms in search of answers. Ancient enemies and unexpected allies await as the threat of Ragnarök grows ever closer.',                                                                                                59.99, 'https://cdn2.steamgriddb.com/grid/fd93fd1de50e084dd7d3b0b9f6950450.png'),
('Monster Hunter: World',   'Welcome to a new world! In Monster Hunter: World, the latest installment in the series, you can enjoy the ultimate hunting experience using everything at your disposal to hunt monsters in a new world teeming with surprises.',                                                                      29.99, 'https://cdn2.steamgriddb.com/grid/0592cd603d4801529227f55223d0d612.png'),
('Pathfinder: Wrath of the Righteous', 'The demon armies have broken through the wardstone barrier and flooded the lands with evil. You are the Commander. Lead crusaders, make impactful choices, and forge your path through an epic fantasy realm.',                                                                           49.99, 'https://cdn2.steamgriddb.com/grid/cf97bd7833d5ae672c9ed3a75bd3744b.jpg'),
('Disco Elysium – The Final Cut', 'You''re a detective with a unique skill system at your disposal and a whole city block to carve your way through. Interrogate unforgettable characters, crack murders, or deal with the impossible pull of the human condition.',                                                              39.99, 'https://cdn2.steamgriddb.com/grid/1a99aae44487990b26c0d98c40afcb49.png'),
('Armored Core VI: Fires of Rubicon', 'Take on the role of an independent mercenary armed with a fully customizable mech, and dive into the world of Rubicon 3. Break through enemy lines and uncover the secrets of a mysterious new substance.',                                                                              59.99, 'https://cdn2.steamgriddb.com/grid/3e9e3298971e5dde4a34d30005508781-fakepng.png'), 
("Tom Clancy's The Division 2",              'Join The Division in this open world tactical shooter. Engage in intense cover-based combat, team up with friends, and build your ultimate agent to save Washington, D.C.',2.99,  'https://cdn2.steamgriddb.com/grid/b1df38b8629b8ced5db21200b342c451-fakepng.png'),
("Tom Clancy's The Division",                'A tactical open world third-person shooter set in New York after a deadly pandemic. Team up or go solo to restore order against deadly factions, secure loot, and conquer high-risk endgame challenges.',                                                                                                                                                                                                                                                                                                  2.99,  'https://cdn2.steamgriddb.com/grid/a1c7facb817c34ea64cc8c2d9c2a6af8-fakepng.png'),
('MONSTER HUNTER WILDS',                     'The unbridled force of nature runs wild and relentless, with environments transforming drastically from one moment to the next. This is a story of monsters and humans and their struggles to live in harmony in a world of duality.',                                                                                                                                                                                                                                                                      69.99, 'https://cdn2.steamgriddb.com/grid/606715151cea04f322521f5c0d7d745b.png'),
('Devil May Cry 5',                          'The ultimate Devil Hunter is back in style, in the game action fans have been waiting for.',                                                                                                                                                                                                                                                                                                                                                                                                             29.99, 'https://cdn2.steamgriddb.com/grid/075db2497f6f4b61e6425a314922f98d.png'),
('Devil May Cry 4 Special Edition',          'The smash-hit stylish action game DMC4 is back and better than ever, boasting new playable characters and game modes!',                                                                                                                                                                                                                                                                                                                                                                                  24.99, 'https://cdn2.steamgriddb.com/grid/0b8fb2585cbd33976511f41030dd75e3.png'),
('Devil May Cry HD Collection',              'The popular stylish action games Devil May Cry, Devil May Cry 2, and Devil May Cry 3 Special Edition return in one collection!',                                                                                                                                                                                                                                                                                                                                                                          29.99, 'https://cdn2.steamgriddb.com/grid/18256dbcecf506006d15b9dc922f1342.png'),
("Devil May Cry 3: Dante's Awakening - SE", "Set before the original Devil May Cry, Devil May Cry 3: Special Edition follows Dante's journey as he battles his evil twin brother in a final confrontation. Along the way, you'll encounter both formidable foes and unusual characters. When attacking, you can choose to use short-range, long-range, defensive, or evasive play styles. Each system has its own set of gameplay features, and you can choose to power up your favorite style for increased special attack abilities. The Special Edition adds the ability to play as Dante's nemesis Virgil and a survival game mode.", 4.79, 'https://cdn2.steamgriddb.com/grid/eeda7b9882543ce81b18415f5e6b42e4.png');

INSERT INTO product_lists (product_id, category_id, user_id, release_date) VALUES 
-- DARK SOULS III (1)
(1,  1, 1, now()),   -- Dark Fantasy
(1,  3, 1, now()),   -- RPG
(1,  4, 1, now()),   -- Adventure

-- DARK SOULS: REMASTERED (2)
(2,  1, 1, now()),   -- Dark Fantasy
(2,  3, 1, now()),   -- RPG
(2,  4, 1, now()),   -- Adventure

-- DARK SOULS II: Scholar of the First Sin (3)
(3,  1, 1, now()),   -- Dark Fantasy
(3,  3, 1, now()),   -- RPG
(3,  4, 1, now()),   -- Adventure

-- ELDEN RING (4)
(4,  1, 1, now()),   -- Dark Fantasy
(4,  3, 1, now()),   -- RPG
(4,  4, 1, now()),   -- Adventure
(4,  9, 1, now()),   -- Open World

-- ELDEN RING NIGHTREIGN (5)
(5,  1, 1, now()),   -- Dark Fantasy
(5,  3, 1, now()),   -- RPG
(5,  4, 1, now()),   -- Adventure

-- ELDEN RING Shadow of the Erdtree (6)
(6,  1, 1, now()),   -- Dark Fantasy
(6,  3, 1, now()),   -- RPG
(6,  4, 1, now()),   -- Adventure

-- Bloodborne (7)
(7,  1, 1, now()),   -- Dark Fantasy
(7,  2, 1, now()),   -- Atmospheric
(7,  8, 1, now()),   -- Horror
(7,  4, 1, now()),   -- Adventure

-- Sekiro: Shadows Die Twice (8)
(8,  1, 1, now()),   -- Dark Fantasy
(8,  2, 1, now()),   -- Atmospheric
(8,  4, 1, now()),   -- Adventure

-- Demon's Souls (9)
(9,  1, 1, now()),   -- Dark Fantasy
(9,  3, 1, now()),   -- RPG
(9,  4, 1, now()),   -- Adventure

-- Lies of P (10)
(10,  1, 1, now()),  -- Dark Fantasy
(10,  2, 1, now()),  -- Atmospheric
(10,  4, 1, now()),  -- Adventure

-- Lords of the Fallen (11)
(11,  1, 1, now()),  -- Dark Fantasy
(11,  3, 1, now()),  -- RPG
(11,  4, 1, now()),  -- Adventure

-- Hollow Knight (12)
(12,  1, 1, now()),  -- Dark Fantasy
(12,  2, 1, now()),  -- Atmospheric
(12, 12, 1, now()),  -- Indie
(12, 13, 1, now()),  -- Platformer

-- Hollow Knight: Silksong (13)
(13,  1, 1, now()),  -- Dark Fantasy
(13,  2, 1, now()),  -- Atmospheric
(13, 12, 1, now()),  -- Indie
(13, 13, 1, now()),  -- Platformer

-- Blasphemous (14)
(14,  1, 1, now()),  -- Dark Fantasy
(14,  2, 1, now()),  -- Atmospheric
(14, 12, 1, now()),  -- Indie
(14, 13, 1, now()),  -- Platformer

-- Blasphemous 2 (15)
(15,  1, 1, now()),  -- Dark Fantasy
(15,  2, 1, now()),  -- Atmospheric
(15, 12, 1, now()),  -- Indie
(15, 13, 1, now()),  -- Platformer

-- Code Vein (16)
(16,  1, 1, now()),  -- Dark Fantasy
(16,  3, 1, now()),  -- RPG
(16,  4, 1, now()),  -- Adventure

-- Mortal Shell (17)
(17,  1, 1, now()),  -- Dark Fantasy
(17,  2, 1, now()),  -- Atmospheric
(17,  8, 1, now()),  -- Horror
(17,  4, 1, now()),  -- Adventure

-- Nioh 2 (18)
(18,  1, 1, now()),  -- Dark Fantasy
(18,  3, 1, now()),  -- RPG
(18,  4, 1, now()),  -- Adventure

-- Wo Long: Fallen Dynasty (19)
(19,  1, 1, now()),  -- Dark Fantasy
(19,  3, 1, now()),  -- RPG
(19,  4, 1, now()),  -- Adventure

-- Ghost of Tsushima (20)
(20,  2, 1, now()),  -- Atmospheric
(20,  3, 1, now()),  -- RPG
(20,  4, 1, now()),  -- Adventure
(20,  9, 1, now()),  -- Open World

-- Dragon's Dogma 2 (21)
(21,  1, 1, now()),  -- Dark Fantasy
(21,  3, 1, now()),  -- RPG
(21,  4, 1, now()),  -- Adventure
(21,  9, 1, now()),  -- Open World

-- The Witcher 3 (22)
(22,  2, 1, now()),  -- Atmospheric
(22,  3, 1, now()),  -- RPG
(22,  4, 1, now()),  -- Adventure
(22,  9, 1, now()),  -- Open World

-- Cyberpunk 2077 (23)
(23,  2, 1, now()),  -- Atmospheric
(23,  3, 1, now()),  -- RPG
(23,  4, 1, now()),  -- Adventure
(23,  9, 1, now()),  -- Open World

-- Baldur's Gate 3 (24)
(24,  1, 1, now()),  -- Dark Fantasy
(24,  3, 1, now()),  -- RPG
(24,  4, 1, now()),  -- Adventure

-- Divinity: Original Sin 2 (25)
(25,  2, 1, now()),  -- Atmospheric
(25,  3, 1, now()),  -- RPG
(25,  4, 1, now()),  -- Adventure

-- Hades (26)
(26,  1, 1, now()),  -- Dark Fantasy
(26,  3, 1, now()),  -- RPG
(26, 12, 1, now()),  -- Indie

-- Hades II (27)
(27,  1, 1, now()),  -- Dark Fantasy
(27,  3, 1, now()),  -- RPG
(27, 12, 1, now()),  -- Indie

-- Dead Cells (28)
(28,  1, 1, now()),  -- Dark Fantasy
(28, 12, 1, now()),  -- Indie
(28, 13, 1, now()),  -- Platformer

-- Salt and Sanctuary (29)
(29,  1, 1, now()),  -- Dark Fantasy
(29, 12, 1, now()),  -- Indie
(29, 13, 1, now()),  -- Platformer

-- Remnant II (30)
(30,  1, 1, now()),  -- Dark Fantasy
(30,  3, 1, now()),  -- RPG
(30,  4, 1, now()),  -- Adventure

-- Star Wars Jedi: Survivor (31)
(31,  2, 1, now()),  -- Atmospheric
(31,  3, 1, now()),  -- RPG
(31,  4, 1, now()),  -- Adventure

-- Final Fantasy XVI (32)
(32,  1, 1, now()),  -- Dark Fantasy
(32,  3, 1, now()),  -- RPG
(32,  4, 1, now()),  -- Adventure

-- God of War (33)
(33,  1, 1, now()),  -- Dark Fantasy
(33,  3, 1, now()),  -- RPG
(33,  4, 1, now()),  -- Adventure

-- God of War Ragnarök (34)
(34,  1, 1, now()),  -- Dark Fantasy
(34,  3, 1, now()),  -- RPG
(34,  4, 1, now()),  -- Adventure

-- Monster Hunter: World (35)
(35,  3, 1, now()),  -- RPG
(35,  4, 1, now()),  -- Adventure
(35,  9, 1, now()),  -- Open World

-- Pathfinder: Wrath of the Righteous (36)
(36,  1, 1, now()),  -- Dark Fantasy
(36,  3, 1, now()),  -- RPG
(36,  4, 1, now()),  -- Adventure

-- Disco Elysium (37)
(37,  2, 1, now()),  -- Atmospheric
(37,  3, 1, now()),  -- RPG
(37, 12, 1, now()),  -- Indie

-- Armored Core VI (38)
(38,  2, 1, now()),  -- Atmospheric
(38,  4, 1, now()),  -- Adventure
(38,  5, 1, now()),  -- Strategy

-- Tom Clancy's The Division 2 (39)
(39,  5, 1, now()),  -- Strategy
(39,  9, 1, now()),  -- Open World
(39, 14, 1, now()),  -- Massively Multiplayer

-- Tom Clancy's The Division (40)
(40,  5, 1, now()),  -- Strategy
(40,  9, 1, now()),  -- Open World
(40, 14, 1, now()),  -- Massively Multiplayer

-- MONSTER HUNTER WILDS (41)
(41,  3, 1, now()),  -- RPG
(41,  4, 1, now()),  -- Adventure
(41,  9, 1, now()),  -- Open World

-- Devil May Cry 5 (42)
(42,  1, 1, now()),  -- Dark Fantasy
(42,  4, 1, now()),  -- Adventure
(42, 15, 1, now()),  -- Fighting

-- Devil May Cry 4 Special Edition (43)
(43,  1, 1, now()),  -- Dark Fantasy
(43,  4, 1, now()),  -- Adventure
(43, 15, 1, now()),  -- Fighting

-- Devil May Cry HD Collection (44)
(44,  1, 1, now()),  -- Dark Fantasy
(44,  4, 1, now()),  -- Adventure
(44, 15, 1, now()),  -- Fighting

-- Devil May Cry 3: Dante's Awakening (45)
(45,  1, 1, now()),  -- Dark Fantasy
(45,  4, 1, now()),  -- Adventure
(45, 15, 1, now());  -- Fighting
