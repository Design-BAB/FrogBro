# Frog Bro

A 2D platformer game built with Go and Raylib where players control a frog collecting flies across challenging vertical platforms.

## About

Frog Bros is a precision platformer featuring:
- **Double-jump mechanics** for navigating multi-tiered platforms
- **Collectible flies** scattered throughout the level
- **Smooth sprite animations** with automatic flipping based on movement direction
- **Physics-based collision detection** using overlap/penetration resolution

## Technical Details

- **Language:** Go
- **Graphics Library:** raylib-go
- **Resolution:** 800x672
- **Target FPS:** 60

## Controls

- **Arrow Keys (Left/Right):** Move the frog
- **Up Arrow:** Jump (press twice for double jump)
- **ESC:** Close the game

## Game Mechanics

The player navigates a series of platforms arranged vertically, collecting animated flies while avoiding falling. The double-jump system allows for precise platforming across gaps and between tiers.

## Code Structure

- Overlap-based collision resolution for smooth block interactions
- Sprite sheet animation system with configurable frame delays
- Modular block placement using row generation helpers
- Separate physics handling for X and Y axes
