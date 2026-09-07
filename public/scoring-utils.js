function performanceLevelFromPct(pct) {
    if (pct > 90) return 'Proficient';
    if (pct >= 70) return 'Competent';
    return 'Not Acceptable';
}

function performanceCardClass(pct) {
    if (pct > 90) return 'score-proficient';
    if (pct >= 70) return 'score-competent';
    return 'score-not-acceptable';
}

function performanceBadgeClass(level) {
    if (level === 'Proficient') return 'bg-success';
    if (level === 'Competent') return 'bg-warning text-dark';
    return 'bg-danger';
}

function scoreResponse(response, weight) {
    if (response === 'Yes') return { points: weight, counts: true };
    if (response === 'No') return { points: 0, counts: true };
    return { points: 0, counts: false };
}

function percentageScore(achieved, possible) {
    if (!possible) return 0;
    return (achieved / possible) * 100;
}
