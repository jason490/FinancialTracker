import { Navigate } from "@solidjs/router";

// TagsRedirect sends legacy /tags bookmarks to the Settings Tags tab.
export default function TagsRedirect() {
  return <Navigate href="/settings?tab=tags" />;
}
