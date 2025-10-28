# PairWise User Manual

## Table of Contents

1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [PairWise Methodology Overview](#p-wvc-methodology-overview)
4. [Using the Application](#using-the-application)
5. [Workflow Guide](#workflow-guide)
6. [Best Practices](#best-practices)
7. [Troubleshooting](#troubleshooting)
8. [Frequently Asked Questions](#frequently-asked-questions)

## Introduction

### What is PairWise?

PairWise (Pairwise-Weighted Value/Complexity) is a structured methodology for objective feature prioritization through group consensus. It combines:

- **Pairwise Comparisons**: Head-to-head feature comparisons for Value and Complexity
- **Win-Count Weighting**: Mathematical calculation of relative feature standings
- **Fibonacci Scoring**: Absolute magnitude scoring using the Fibonacci sequence
- **Consensus Requirement**: All scores must be agreed upon by the entire team

### Benefits of PairWise

- **Objective Decision Making**: Removes subjective bias from prioritization
- **Team Alignment**: Ensures all stakeholders agree on priorities
- **Transparent Process**: Clear methodology that everyone can understand
- **Mathematical Rigor**: Uses proven statistical methods for ranking
- **Collaborative**: Encourages team discussion and consensus building

### When to Use PairWise

PairWise is ideal for:

- Product feature prioritization
- Project planning and resource allocation
- Strategic initiative ranking
- Requirements prioritization
- Any scenario requiring group consensus on competing options

## Getting Started

### System Requirements

- **Web Browser**: Chrome, Firefox, Safari, or Edge (latest versions)
- **Screen Resolution**: Minimum 1024x768, recommended 1920x1080
- **Internet Connection**: Required for real-time collaboration
- **Team Size**: 3-12 participants (optimal: 5-7)

### Accessing the Application

1. Open your web browser
2. Navigate to the PairWise application URL provided by your administrator
3. The application loads directly in your browser - no installation required

### User Roles

- **Facilitator**: Manages the session, guides the process, resolves conflicts
- **Participant**: Votes in comparisons, provides Fibonacci scores, participates in discussions
- **Observer**: Can view the process but cannot vote (optional role)

## PairWise Methodology Overview

### The Four-Step Process

#### Step 1: Pairwise Comparisons (Value)

- Compare each feature pair for **business value**
- Ask: "Which feature provides more value to users/business?"
- Vote as a team until consensus is reached
- Results in win-count weights for value (WValue)

#### Step 2: Pairwise Comparisons (Complexity)

- Compare each feature pair for **implementation complexity**
- Ask: "Which feature is more complex to implement?"
- Vote as a team until consensus is reached
- Results in win-count weights for complexity (WComplexity)

#### Step 3: Fibonacci Scoring (Value & Complexity)

- Assign absolute magnitude scores using Fibonacci numbers (1, 2, 3, 5, 8, 13, 21, 34, 55, 89)
- Score each feature for both value (SValue) and complexity (SComplexity)
- Team must reach consensus on all scores

#### Step 4: Final Priority Calculation

- Calculate Final Priority Score (FPS) using the formula:
  ```
  FPS = (SValue × WValue) / (SComplexity × WComplexity)
  ```
- Features are ranked by FPS (highest = highest priority)

### Mathematical Foundation

The PairWise formula balances value and complexity:

- **Numerator (SValue × WValue)**: Weighted value benefit
- **Denominator (SComplexity × WComplexity)**: Weighted implementation cost
- **Result**: Value-to-complexity ratio for objective prioritization

## Using the Application

### Creating a New Project

1. **Start a New Project**

   - Click "New Project" on the homepage
   - Enter project name and description
   - Click "Create Project"

2. **Project Setup**
   - Review project details
   - Proceed to attendee management

### Managing Attendees

1. **Add Team Members**

   - Click "Add Attendee"
   - Enter name and role
   - Designate one person as facilitator
   - Add all team members (minimum 2 required)

2. **Facilitator Responsibilities**
   - Guide the session process
   - Resolve voting conflicts
   - Ensure consensus is reached
   - Move between workflow phases

### Adding Features

1. **Feature Entry**

   - Click "Add Feature"
   - Enter feature title (be concise and clear)
   - Provide detailed description
   - Add acceptance criteria (optional but recommended)

2. **Feature Requirements**

   - Minimum 2 features required for PairWise
   - Recommended: 5-15 features for optimal comparison
   - Ensure features are similar in scope for fair comparison

3. **Bulk Import (Optional)**
   - Use "Import Features" for large feature sets
   - Download CSV template
   - Fill in feature details
   - Upload completed CSV file

### Pairwise Comparison Process

#### Value Comparisons

1. **Starting the Session**

   - Navigate to "Pairwise Comparisons"
   - Select "Value" criterion
   - Click "Start Session"

2. **Making Comparisons**

   - Review each feature pair presented
   - Discuss with your team which feature provides more business value
   - Consider factors like:
     - User benefit and satisfaction
     - Revenue potential
     - Strategic alignment
     - Market competitive advantage
     - Risk reduction

3. **Voting Process**

   - Each team member votes independently
   - Choose the feature with higher value OR select "Tie"
   - Results are shown in real-time
   - If consensus is not reached, discuss and vote again

4. **Reaching Consensus**
   - All team members must agree on the outcome
   - Facilitator guides discussion if there are conflicts
   - Consider all perspectives before finalizing votes

#### Complexity Comparisons

1. **Starting Complexity Session**

   - Select "Complexity" criterion after value comparisons
   - Click "Start Session"

2. **Complexity Factors to Consider**

   - Technical difficulty
   - Resource requirements
   - Time to implement
   - Dependencies on other systems
   - Risk and uncertainty
   - Testing complexity

3. **Voting for Complexity**
   - Vote for the feature that is MORE complex to implement
   - Higher complexity = more difficult/expensive to build
   - Reach consensus on each comparison

### Fibonacci Scoring

#### Understanding Fibonacci Numbers

The Fibonacci sequence (1, 2, 3, 5, 8, 13, 21, 34, 55, 89) provides natural scaling:

- **1**: Minimal impact/effort
- **2**: Very low impact/effort
- **3**: Low impact/effort
- **5**: Medium impact/effort
- **8**: High impact/effort
- **13**: Very high impact/effort
- **21**: Extremely high impact/effort
- **34+**: Massive impact/effort (consider breaking down)

#### Value Scoring Process

1. **Access Fibonacci Scoring**

   - Navigate to "Fibonacci Scoring"
   - Select "Value" tab

2. **Scoring Guidelines**

   - Consider the absolute magnitude of business value
   - Compare against the Fibonacci scale
   - Think about user impact, revenue potential, strategic value
   - Discuss as a team and reach consensus

3. **Example Value Scores**
   - **1-2**: Nice-to-have features, minor improvements
   - **3-5**: Important features with clear user benefit
   - **8-13**: Major features with significant business impact
   - **21+**: Game-changing features, core differentiators

#### Complexity Scoring Process

1. **Select Complexity Tab**

   - Navigate to "Complexity" in Fibonacci scoring

2. **Complexity Assessment**

   - Consider implementation difficulty
   - Account for technical debt and dependencies
   - Factor in team expertise and resources
   - Include testing and deployment complexity

3. **Example Complexity Scores**
   - **1-2**: Simple changes, configuration updates
   - **3-5**: Standard development work, known patterns
   - **8-13**: Complex features requiring significant effort
   - **21+**: Massive undertakings, major architectural changes

### Viewing Results

1. **Priority Rankings**

   - Navigate to "Results" after completing all phases
   - View features ranked by Final Priority Score (FPS)
   - Higher FPS = higher priority

2. **Understanding the Results**

   - **Rank**: Overall priority position
   - **Feature**: Name and description
   - **FPS**: Final Priority Score (value/complexity ratio)
   - **Value Metrics**: SValue × WValue
   - **Complexity Metrics**: SComplexity × WComplexity

3. **Exporting Results**
   - Click "Export Results" for CSV download
   - Share with stakeholders and development teams
   - Use for sprint planning and roadmap creation

## Workflow Guide

### Phase Navigation

The application guides you through eight distinct phases:

#### 1. Setup Phase ✅

- Create project and define scope
- Review PairWise methodology with team
- **Completion Criteria**: Project created with clear objectives

#### 2. Attendees Phase 👥

- Add all team members
- Assign facilitator role
- **Completion Criteria**: Minimum 2 attendees, 1 facilitator assigned

#### 3. Features Phase 📋

- Add all features to be prioritized
- Ensure feature descriptions are clear
- **Completion Criteria**: Minimum 2 features added

#### 4. Pairwise Value Phase 💰

- Compare all feature pairs for business value
- Reach consensus on each comparison
- **Completion Criteria**: All value comparisons completed with consensus

#### 5. Pairwise Complexity Phase ⚙️

- Compare all feature pairs for implementation complexity
- Reach consensus on each comparison
- **Completion Criteria**: All complexity comparisons completed with consensus

#### 6. Fibonacci Value Phase 📊

- Score each feature's absolute value magnitude
- Use Fibonacci numbers (1, 2, 3, 5, 8, 13, 21, 34, 55, 89)
- **Completion Criteria**: All features have consensus value scores

#### 7. Fibonacci Complexity Phase 🔧

- Score each feature's absolute complexity magnitude
- Use Fibonacci numbers for consistency
- **Completion Criteria**: All features have consensus complexity scores

#### 8. Results Phase 🏆

- Calculate Final Priority Scores
- Review rankings and discuss implications
- **Completion Criteria**: Priority calculations complete, results reviewed

### Session Management Tips

#### For Facilitators

1. **Preparation**

   - Review all features before starting
   - Ensure team understands PairWise methodology
   - Allocate sufficient time (2-4 hours for 10 features)

2. **During the Session**

   - Keep discussions focused and time-boxed
   - Encourage all participants to voice opinions
   - Guide team to consensus without forcing agreement
   - Take breaks as needed

3. **Managing Conflicts**
   - Listen to all perspectives
   - Focus on objective criteria
   - Use data and examples to support discussions
   - Remind team of scoring guidelines

#### For Participants

1. **Come Prepared**

   - Review features before the session
   - Understand your organization's priorities
   - Bring relevant domain expertise

2. **Active Participation**
   - Voice your opinions clearly
   - Listen to other perspectives
   - Ask questions when uncertain
   - Support consensus building

### Real-Time Collaboration

#### Live Updates

- All participants see votes and changes immediately
- Progress indicators show session status
- Notifications alert when consensus is reached

#### Communication Features

- Use built-in chat or external communication tools
- Share screens for detailed feature discussions
- Document decisions for future reference

## Best Practices

### Pre-Session Preparation

#### Define Clear Scope

- Ensure all features are similar in scope
- Break down large features into smaller, comparable pieces
- Remove obvious non-starters before the session

#### Team Preparation

- Share feature list with team in advance
- Provide context and background information
- Ensure key stakeholders participate

#### Time Management

- Allocate 15-30 minutes per feature for pairwise comparisons
- Plan for breaks every 1-2 hours
- Consider splitting large sessions across multiple days

### During the Session

#### Maintain Focus

- Keep discussions centered on the specific criterion (value or complexity)
- Avoid getting sidetracked by implementation details during value comparisons
- Stay objective and data-driven

#### Encourage Participation

- Ensure all voices are heard
- Rotate who speaks first on comparisons
- Address conflicts constructively

#### Documentation

- Record key insights and assumptions
- Note any features that need further research
- Capture action items for follow-up

### Post-Session Actions

#### Results Review

- Validate that rankings make intuitive sense
- Identify any surprising results for further discussion
- Consider if any features need to be re-evaluated

#### Implementation Planning

- Use rankings to inform sprint planning
- Consider dependencies when ordering development
- Communicate results to broader stakeholder group

## Troubleshooting

### Common Issues and Solutions

#### Session Won't Start

**Problem**: Cannot begin pairwise comparisons
**Solutions**:

- Ensure minimum 2 attendees are added
- Verify at least 2 features exist
- Check that all required fields are completed
- Refresh browser if interface seems stuck

#### Votes Not Registering

**Problem**: Clicking vote buttons doesn't record votes
**Solutions**:

- Check internet connection
- Refresh browser page
- Clear browser cache and cookies
- Try different browser or incognito mode
- Verify you're logged in as the correct attendee

#### Consensus Not Reached

**Problem**: Team cannot agree on comparisons
**Solutions**:

- Revisit feature definitions for clarity
- Focus discussion on specific criterion (value vs complexity)
- Consider if features need to be redefined or split
- Take a break and return with fresh perspective
- Have facilitator guide structured discussion

#### Results Seem Wrong

**Problem**: Final rankings don't match expectations
**Solutions**:

- Review pairwise comparison results for accuracy
- Verify Fibonacci scores make sense
- Check if any votes were entered incorrectly
- Consider team biases in original scoring
- Discuss specific ranking surprises with team

#### Technical Issues

**Problem**: Application crashes or freezes
**Solutions**:

- Refresh browser and reload application
- Check system requirements and browser compatibility
- Clear browser cache and restart
- Contact system administrator for server issues
- Try accessing from different device or network

### Getting Help

#### In-Application Support

- Use built-in help tooltips and guidance
- Check workflow progress indicators
- Review error messages for specific guidance

#### Team Support

- Consult with facilitator for process questions
- Engage team members for feature clarification
- Escalate to project sponsor for scope issues

#### Technical Support

- Contact system administrator for technical issues
- Provide specific error messages and browser information
- Include steps to reproduce the problem

## Frequently Asked Questions

### General PairWise Questions

**Q: How many features should we include in a PairWise session?**
A: Optimal range is 5-15 features. Fewer than 5 provides limited comparison value; more than 15 becomes time-intensive. For larger feature sets, consider grouping into themes or conducting multiple sessions.

**Q: How long does a PairWise session take?**
A: Typically 2-4 hours for 8-10 features, depending on team size and feature complexity. Plan 15-30 minutes per feature for thorough discussion.

**Q: What if our team can't reach consensus?**
A: Focus discussions on objective criteria, ensure all perspectives are heard, and consider if features need clearer definitions. The facilitator should guide structured discussion without forcing agreement.

**Q: Can we modify features during the session?**
A: Yes, but be cautious. Small clarifications are fine, but major changes may invalidate previous comparisons. Consider pausing to update feature descriptions if needed.

### Process Questions

**Q: Why do we compare pairs instead of ranking all features at once?**
A: Pairwise comparison is more reliable because human brains are better at comparing two options than ranking many. It also reduces position bias and anchoring effects.

**Q: What's the difference between win-count weights and Fibonacci scores?**
A: Win-count weights (W) show relative standings from pairwise comparisons. Fibonacci scores (S) show absolute magnitude. Combined, they provide both relative ranking and absolute scale.

**Q: Can the same person be both facilitator and participant?**
A: Yes, the facilitator can vote, but they should prioritize guiding the process and ensuring all voices are heard over advocating for specific positions.

**Q: What if we discover new features during the session?**
A: You can add features, but this requires redoing comparisons involving the new feature. Consider if the new feature is truly essential or can wait for a future session.

### Technical Questions

**Q: Can we use PairWise remotely?**
A: Yes, the application is designed for remote collaboration with real-time updates and vote sharing. Combine with video conferencing for best results.

**Q: Is our data secure?**
A: Check with your system administrator about data security policies. The application typically stores data securely, but specific security measures depend on your deployment.

**Q: Can we export our results?**
A: Yes, results can be exported to CSV format for use in other tools like Excel, project management software, or presentation tools.

**Q: What browsers are supported?**
A: Modern versions of Chrome, Firefox, Safari, and Edge. For best performance, use the latest version of your preferred browser.

### Results and Follow-up

**Q: How should we use the PairWise results?**
A: Use rankings to inform development priorities, sprint planning, and resource allocation. Remember that PairWise provides one input to decision-making, not the final word.

**Q: Can we re-run PairWise for the same features?**
A: Yes, priorities may change over time due to market conditions, technical discoveries, or strategic shifts. Consider re-running quarterly or when significant changes occur.

**Q: What if development reveals our complexity estimates were wrong?**
A: This is normal and expected. Use actual complexity data to refine future estimates and consider adjusting priorities based on new information.

**Q: Should we prioritize features with the highest FPS scores?**
A: Generally yes, but also consider dependencies, strategic timing, and resource availability. PairWise provides prioritization guidance within broader strategic context.

---

## Getting Started Checklist

Use this checklist to ensure successful PairWise sessions:

### Pre-Session Preparation

- [ ] Project scope and objectives clearly defined
- [ ] All key stakeholders identified and available
- [ ] Feature list prepared and reviewed
- [ ] Time allocated (2-4 hours for typical session)
- [ ] Technology setup tested (browsers, connectivity)

### Team Setup

- [ ] Facilitator designated and briefed on process
- [ ] All attendees added to project
- [ ] Roles and responsibilities clarified
- [ ] PairWise methodology explained to team

### Feature Preparation

- [ ] Minimum 2 features defined (recommend 5-15)
- [ ] Feature descriptions are clear and comparable
- [ ] Acceptance criteria provided where helpful
- [ ] Features are similar in scope and granularity

### Session Execution

- [ ] Pairwise value comparisons completed
- [ ] Pairwise complexity comparisons completed
- [ ] Fibonacci value scores assigned
- [ ] Fibonacci complexity scores assigned
- [ ] Consensus reached on all scoring decisions

### Results and Follow-up

- [ ] Final priority rankings reviewed
- [ ] Results exported and shared with stakeholders
- [ ] Next steps and action items documented
- [ ] Development priorities updated based on results

---

## Additional Resources

### Training Materials

- PairWise methodology overview presentation
- Video tutorials for each workflow phase
- Practice exercises with sample data

### Templates and Tools

- Feature definition template
- Session agenda template
- Results presentation template
- CSV import/export formats

### Support Contacts

- Technical support: [Contact information]
- Process questions: [Facilitator resources]
- Feature requests: [Product feedback]

Remember: PairWise is a powerful tool for objective prioritization, but it's most effective when combined with good communication, clear feature definitions, and commitment to consensus-building within your team.
