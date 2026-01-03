/**
 * @param {import('@actions/github').GitHub} github
 * @param {import('@actions/github').context} context
 */
module.exports = async ({ github, context }) => {
  console.log("=== ASSIGN BOT START ===");

  const comment = context.payload.comment.body.trim();
  const issue = context.payload.issue;
  const commenter = context.payload.comment.user.login;
  const owner = context.repo.owner;
  const repo = context.repo.repo;
  const issue_number = issue.number;

  console.log("Comment:", `'${comment}'`);
  console.log("Issue #", issue_number);
  console.log("Commenter:", commenter);

  // Ignore bots FIRST
  if (context.payload.comment.user.type === "Bot") {
    console.log("Bot ignored");
    return;
  }

  const assignees = issue.assignees.map((a) => a.login);
  const labels = issue.labels.map((l) => (typeof l === "string" ? l : l.name));

  console.log("Assignees:", assignees);
  console.log("Labels:", labels);

  const APPROVAL_LABEL = "approved";
  const isApproved = labels.includes(APPROVAL_LABEL);
  console.log("Approved?", isApproved);

  // /assign
  if (comment === "/assign") {
    console.log("/assign");
    if (!isApproved) {
      console.log("No approved");
      await github.rest.issues.createComment({
        owner,
        repo,
        issue_number,
        body: `This issue needs \`${APPROVAL_LABEL}\` label first.`,
      });
      return;
    }
    if (assignees.length === 0) {
      console.log("Assigning", commenter);
      await github.rest.issues.addAssignees({
        owner,
        repo,
        issue_number,
        assignees: [commenter],
      });
      await github.rest.issues.createComment({
        owner,
        repo,
        issue_number,
        body: `@${commenter} assigned!`,
      });
    } else {
      console.log("Already:", assignees);
      await github.rest.issues.createComment({
        owner,
        repo,
        issue_number,
        body: `Already @${assignees.join(", @")}. Use /unassign.`,
      });
    }
    return;
  }

  // /unassign
  if (comment === "/unassign") {
    console.log("/unassign");
    if (assignees.includes(commenter)) {
      console.log("Unassigning");
      await github.rest.issues.removeAssignees({
        owner,
        repo,
        issue_number,
        assignees: [commenter],
      });
      await github.rest.issues.createComment({
        owner,
        repo,
        issue_number,
        body: `@${commenter} unassigned.`,
      });
    } else {
      console.log("Not assigned");
      await github.rest.issues.createComment({
        owner,
        repo,
        issue_number,
        body: `Only assignee can /unassign.`,
      });
    }
    return;
  }

  console.log("No command");
  console.log("=== ASSIGN BOT END ===");
};
