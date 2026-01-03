/**
 * @param {import('@actions/github').GitHub} github
 * @param {import('@actions/github').context} context
 */

const comment = context.payload.comment.body.trim();
const issue = context.payload.issue;
const commenter = context.payload.comment.user.login;

const owner = context.repo.owner;
const repo = context.repo.repo;
const issue_number = issue.number;

// Ignore bots
if (context.payload.comment.user.type === "Bot") {
  return;
}

const assignees = issue.assignees.map(a => a.login);
const labels = issue.labels.map(l =>
  typeof l === "string" ? l : l.name
);

const APPROVAL_LABEL = "approved";
const isApproved = labels.includes(APPROVAL_LABEL);

// /assign command
if (comment === "/assign") {
  if (!isApproved) {
    await github.rest.issues.createComment({
      owner,
      repo,
      issue_number,
      body: `⛔ This issue is not approved yet.\n\nA maintainer must add the \`${APPROVAL_LABEL}\` label before assignment.`,
    });
    return;
  }

  if (assignees.length === 0) {
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
      body: `✅ @${commenter} has been assigned to this issue.`,
    });
  } else {
    await github.rest.issues.createComment({
      owner,
      repo,
      issue_number,
      body: `⚠️ This issue is already assigned to @${assignees.join(
        ", @"
      )}.\n\nIf you are no longer working on it, please comment \`/unassign\`.`,
    });
  }
}

// /unassign command
if (comment === "/unassign") {
  if (assignees.includes(commenter)) {
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
      body: `🔓 @${commenter} has unassigned themselves. The issue is now available.`,
    });
  } else {
    await github.rest.issues.createComment({
      owner,
      repo,
      issue_number,
      body: `❌ Only the currently assigned user can unassign themselves.`,
    });
  }
}
