#!/bin/bash

nvm install lts/*
npx semantic-release --branches ["main", "feature/*"]