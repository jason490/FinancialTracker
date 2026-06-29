import { clientApiRequest } from "./api";
import type {
  CreateCategoryRequest,
  CreateTagRequest,
  DeleteCategoryRequest,
  MoveTagRequest,
  TagFilterView,
  TagsPayload,
  UpdateCategoryRequest,
  UpdateTagRequest,
} from "./types";

// getTags loads all categories and their tags.
export async function getTags(): Promise<TagsPayload> {
  return clientApiRequest<TagsPayload>("/api/v1/tags");
}

// getTagFilters loads auto-tagging filters for a tag.
export async function getTagFilters(tagId: number): Promise<TagFilterView[]> {
  return clientApiRequest<TagFilterView[]>(`/api/v1/tags/${tagId}/filters`);
}

// createTag creates a new tag and returns the refreshed tags payload.
export async function createTag(body: CreateTagRequest): Promise<TagsPayload> {
  return clientApiRequest<TagsPayload>("/api/v1/tags", { method: "POST", body });
}

// updateTag updates a tag and returns the refreshed tags payload.
export async function updateTag(tagId: number, body: UpdateTagRequest): Promise<TagsPayload> {
  return clientApiRequest<TagsPayload>(`/api/v1/tags/${tagId}`, { method: "PUT", body });
}

// deleteTag removes a tag and returns the refreshed tags payload.
export async function deleteTag(tagId: number): Promise<TagsPayload> {
  return clientApiRequest<TagsPayload>(`/api/v1/tags/${tagId}`, { method: "DELETE" });
}

// moveTag moves a tag to another category and returns the refreshed tags payload.
export async function moveTag(tagId: number, body: MoveTagRequest): Promise<TagsPayload> {
  return clientApiRequest<TagsPayload>(`/api/v1/tags/${tagId}/move`, { method: "POST", body });
}

// createCategory creates a category and returns the refreshed tags payload.
export async function createCategory(body: CreateCategoryRequest): Promise<TagsPayload> {
  return clientApiRequest<TagsPayload>("/api/v1/categories", { method: "POST", body });
}

// updateCategory renames a category and returns the refreshed tags payload.
export async function updateCategory(
  categoryId: number,
  body: UpdateCategoryRequest
): Promise<TagsPayload> {
  return clientApiRequest<TagsPayload>(`/api/v1/categories/${categoryId}`, {
    method: "PUT",
    body,
  });
}

// deleteCategory deletes a category and returns the refreshed tags payload.
export async function deleteCategory(
  categoryId: number,
  body: DeleteCategoryRequest
): Promise<TagsPayload> {
  return clientApiRequest<TagsPayload>(`/api/v1/categories/${categoryId}`, {
    method: "DELETE",
    body,
  });
}
