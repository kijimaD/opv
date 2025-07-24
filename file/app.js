let syncInterval = null;
let animationFrameId = null;
let remainingTime = 0;
let totalTime = 0;
let taskTitle = '';
let todayPoints = 0;
let isActive = false;
let wasActive = false; // Track previous state
let lastSyncTime = 0; // Server time at last sync
let currentData = null; // Store full data from server
let lastUpdateTime = 0; // For animation frame timing

function formatTime(seconds) {
    const minutes = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
}

function showCompletionAnimation() {
    const overlay = document.getElementById('completionOverlay');
    overlay.classList.add('show');

    // Hide after 3 seconds
    setTimeout(() => {
        overlay.classList.remove('show');
    }, 3000);
}

function updateTomatoDisplay(points) {
    const tomatoContainer = document.getElementById('tomatoes');
    const pointsText = document.querySelector('.points-number');

    // Update points text
    pointsText.textContent = `Today: ${points} pomodoros`;

    // Clear existing tomatoes
    tomatoContainer.innerHTML = '';

    // Add tomato SVGs with separators every 4 items
    for (let i = 0; i < points; i++) {
        // Add separator before every group of 4 (except the first group)
        if (i > 0 && i % 4 === 0) {
            const separator = document.createElement('div');
            separator.className = 'separator';
            tomatoContainer.appendChild(separator);
        }

        const tomato = document.createElement('img');
        tomato.className = 'tomato';
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

function updateDisplay(timestamp) {
    if (!lastUpdateTime) lastUpdateTime = timestamp;

    // Update immediately on first call or every second
    if (timestamp === 0 || timestamp - lastUpdateTime >= 1000) {
        lastUpdateTime = timestamp;

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
            container.classList.remove('inactive');
            titleEl.textContent = taskTitle || 'No task';
            timeEl.textContent = formatTime(remainingTime);

            const progress = ((totalTime - remainingTime) / totalTime) * 100;
            progressEl.style.width = Math.min(100, progress) + '%';
        } else {
            document.body.classList.remove('active');
            container.classList.add('inactive');
            titleEl.textContent = 'Pomodoro not active';
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

    // Continue animation loop if active
    if (isActive) {
        animationFrameId = requestAnimationFrame(updateDisplay);
    }
}

async function fetchPomodoroData() {
    try {
        const response = await fetch('/api/pomodoro');
        if (!response.ok) throw new Error('Network response was not ok');

        const data = await response.json();

        // Update global state
        wasActive = isActive; // Store previous state
        currentData = data; // Store full data
        lastSyncTime = data.serverTime || Math.floor(Date.now() / 1000); // Store server time

        remainingTime = data.remainingTime || 0;
        totalTime = data.totalTime || 0;
        taskTitle = data.taskTitle || '';
        todayPoints = data.todayPoints || 0;
        isActive = data.isActive || false;

        // Reset timing for animation frame
        lastUpdateTime = 0;

        // Always update display immediately
        updateDisplay(0);

        // Start/stop animation based on active state
        if (isActive && !animationFrameId) {
            animationFrameId = requestAnimationFrame(updateDisplay);
        } else if (!isActive && animationFrameId) {
            cancelAnimationFrame(animationFrameId);
            animationFrameId = null;
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
        if (syncInterval) clearInterval(syncInterval);
        if (animationFrameId) {
            cancelAnimationFrame(animationFrameId);
            animationFrameId = null;
        }
    } else {
        // When tab becomes visible again
        fetchPomodoroData();
        syncInterval = setInterval(fetchPomodoroData, 10000);
        // Animation will restart via fetchPomodoroData if active
    }
});