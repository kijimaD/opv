let syncInterval = null;
let updateInterval = null; // For 1-second updates
let remainingTime = 0;
let totalTime = 0;
let taskTitle = '';
let todayPoints = 0;
let isActive = false;
let wasActive = false; // Track previous state
let lastSyncTime = 0; // Server time at last sync
let currentData = null; // Store full data from server

function formatTime(seconds) {
    const minutes = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
}

function showCompletionAnimation() {
    const overlay = document.getElementById('completionOverlay');
    overlay.classList.remove('opacity-0', 'invisible');
    overlay.classList.add('opacity-100', 'visible');

    // Hide after 3 seconds
    setTimeout(() => {
        overlay.classList.remove('opacity-100', 'visible');
        overlay.classList.add('opacity-0', 'invisible');
    }, 3000);
}

let lastDisplayedPoints = -1; // Track last displayed points

function updateTomatoDisplay(points) {
    // Skip if points haven't changed
    if (points === lastDisplayedPoints) return;
    lastDisplayedPoints = points;

    const tomatoContainer = document.getElementById('tomatoes');
    const pointsText = document.querySelector('.points-number');

    // Update points text
    pointsText.textContent = `Today: ${points} pmds`;

    // Clear existing tomatoes
    tomatoContainer.innerHTML = '';

    // Add tomato SVGs with separators every 4 items
    for (let i = 0; i < points; i++) {
        // Add separator before every group of 4 (except the first group)
        if (i > 0 && i % 4 === 0) {
            const separator = document.createElement('div');
            separator.className = 'vr mx-2';
            separator.style.height = '32px';
            tomatoContainer.appendChild(separator);
        }

        const tomato = document.createElement('img');
        tomato.className = 'img-fluid';
        tomato.style.width = '32px';
        tomato.style.height = '32px';
        tomato.src = '/file/tomato.svg';
        tomato.alt = 'Pomodoro';
        tomatoContainer.appendChild(tomato);
    }
}

function calculateRemainingTime() {
    if (!currentData || !currentData.isActive) return 0;

    const now = Math.floor(Date.now() / 1000); // Current time in seconds
    const elapsedSinceSync = now - lastSyncTime;
    const remaining = currentData.remainingTime - elapsedSinceSync;

    return Math.max(0, remaining);
}

function updateDisplay() {
    const container = document.querySelector('.container');
    const titleEl = document.getElementById('taskTitle');
    const timeEl = document.getElementById('timeDisplay');
    const progressEl = document.getElementById('progressBar');
    const pointsEl = document.getElementById('todayPoints');

    // Calculate current remaining time based on server sync
    if (isActive) {
        remainingTime = calculateRemainingTime();
    }

    // Update display based on active state
    if (isActive && remainingTime > 0) {
        document.body.classList.add('active');
        titleEl.textContent = taskTitle || 'No task';
        timeEl.textContent = formatTime(remainingTime);

        const progress = ((totalTime - remainingTime) / totalTime) * 100;
        progressEl.style.width = Math.min(100, progress) + '%';
    } else {
        document.body.classList.remove('active');
        titleEl.textContent = 'Are you ready?';
        timeEl.textContent = '--:--';
        progressEl.style.width = '0%';

        // Check if we just completed a pomodoro
        if (remainingTime <= 0 && wasActive && !isActive) {
            showCompletionAnimation();
            wasActive = false; // Prevent multiple animations
        }
    }

    updateTomatoDisplay(todayPoints);
}

async function fetchPomodoroData() {
    try {
        const response = await fetch('/api/pomodoro');
        if (!response.ok) throw new Error('Network response was not ok');

        const data = await response.json();

        // Update global state
        wasActive = isActive; // Store previous state
        currentData = data; // Store full data
        lastSyncTime = Math.floor(Date.now() / 1000); // Store server time

        remainingTime = data.remainingTime || 0;
        totalTime = data.totalTime || 0;
        taskTitle = data.taskTitle || '';
        todayPoints = data.todayPoints || 0;
        isActive = data.isActive || false;

        // Always update display immediately
        updateDisplay();

        // Start/stop 1-second update interval based on active state
        if (isActive && !updateInterval) {
            updateInterval = setInterval(updateDisplay, 1000);
        } else if (!isActive && updateInterval) {
            clearInterval(updateInterval);
            updateInterval = null;
        }
    } catch (error) {
        console.error('Error fetching pomodoro data:', error);
        document.getElementById('taskTitle').textContent = 'Error loading data';
    }
}

// Initialize
window.addEventListener('load', () => {
    fetchPomodoroData();
    syncInterval = setInterval(fetchPomodoroData, 10000);
});

// Handle page visibility
document.addEventListener('visibilitychange', () => {
    if (document.hidden) {
        if (syncInterval) {
            clearInterval(syncInterval);
            syncInterval = null;
        }
        if (updateInterval) {
            clearInterval(updateInterval);
            updateInterval = null;
        }
    } else {
        // When tab becomes visible again
        fetchPomodoroData();
        if (!syncInterval) {
            syncInterval = setInterval(fetchPomodoroData, 10000);
        }
    }
});
