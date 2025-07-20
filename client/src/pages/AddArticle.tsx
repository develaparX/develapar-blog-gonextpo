import React, { useState, useEffect } from "react";
import { useForm, type SubmitHandler } from "react-hook-form";
import { TagsInput } from "./TagsInput";
import { slugify } from "../utils/slugify";

// import QuillEditor from "./QuillEditor";
// import CKEditor5Component from "./CKEditor5Component";

// Interface untuk form data
interface ArticleFormData {
  title: string;
  slug: string;
  content: string;
  category_id: string;
}

// Interface untuk category
interface Category {
  id: number;
  name: string;
}

// Interface untuk article data yang akan dikirim
interface ArticleData extends ArticleFormData {
  tags: string[];
}

const AddArticle: React.FC = () => {
  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors },
  } = useForm<ArticleFormData>();

  const [content, setContent] = useState<string>("");
  const [tags, setTags] = useState<string[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false);

  const title = watch("title", "");

  // Auto-generate slug from title
  useEffect(() => {
    if (title) {
      setValue("slug", slugify(title));
    }
  }, [title, setValue]);

  // Fetch categories
  useEffect(() => {
    setCategories([
      { id: 1, name: "Technology" },
      { id: 2, name: "Business" },
      { id: 3, name: "Health" },
      { id: 4, name: "Education" },
    ]);
  }, []);

  const onSubmit: SubmitHandler<ArticleFormData> = async (data) => {
    setIsLoading(true);
    try {
      const articleData: ArticleData = {
        ...data,
        tags: tags,
      };

      console.log("Article data to submit:", articleData);
      await new Promise((resolve) => setTimeout(resolve, 1000));
      alert("Article published successfully!");
    } catch (error) {
      console.error("Error publishing article:", error);
      alert("Failed to publish article. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  // Handler untuk content change dari WYSIWYG editor
  // const handleContentChange = (newContent: string) => {
  //   setContent(newContent);
  //   setValue("content", newContent);
  // };

  return (
    <div className="w-full mx-auto p-6 bg-white rounded-lg shadow-md">
      <h1 className="text-3xl font-bold text-gray-800 mb-6">
        Create New Article
      </h1>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        {/* Title */}
        <div className="space-y-2">
          <label
            htmlFor="title"
            className="block text-sm font-medium text-gray-700"
          >
            Title
          </label>
          <input
            id="title"
            type="text"
            {...register("title", { required: "Title is required" })}
            className="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            placeholder="Enter article title"
          />
          {errors.title && (
            <p className="text-red-500 text-sm">{errors.title.message}</p>
          )}
        </div>

        {/* Slug */}
        <div className="space-y-2">
          <label
            htmlFor="slug"
            className="block text-sm font-medium text-gray-700"
          >
            Slug
          </label>
          <input
            id="slug"
            type="text"
            {...register("slug", { required: "Slug is required" })}
            className="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            placeholder="article-slug"
          />
          {errors.slug && (
            <p className="text-red-500 text-sm">{errors.slug.message}</p>
          )}
        </div>

        {/* Category */}
        <div className="space-y-2">
          <label
            htmlFor="category_id"
            className="block text-sm font-medium text-gray-700"
          >
            Category
          </label>
          <select
            id="category_id"
            {...register("category_id")}
            className="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          >
            <option value="">Select a category</option>
            {categories.map((category: Category) => (
              <option key={category.id} value={category.id}>
                {category.name}
              </option>
            ))}
          </select>
        </div>

        {/* Tags */}
        <div className="space-y-2">
          <label className="block text-sm font-medium text-gray-700">
            Tags
          </label>
          <TagsInput
            tags={tags}
            setTags={setTags}
            placeholder="Add tags and press Enter"
            maxTags={5}
          />
          <p className="text-sm text-gray-500">
            Press Enter to add a tag. Maximum 5 tags.
          </p>
        </div>

        {/* Content dengan WYSIWYG Editor */}
        <div className="space-y-2">
          <label className="block text-sm font-medium text-gray-700">
            Content
          </label>

          {/* TinyMCE Editor */}

          {/* Alternatif dengan Quill Editor */}
          {/* <QuillEditor
            value={content}
            onChange={handleContentChange}
            height="500px"
          /> */}

          {/* Alternatif dengan CKEditor 5 */}
          {/* <CKEditor5Component
            value={content}
            onChange={handleContentChange}
          /> */}

          {/* Hidden input untuk react-hook-form */}
          <input
            type="hidden"
            {...register("content", { required: "Content is required" })}
            value={content}
          />

          {errors.content && (
            <p className="text-red-500 text-sm">{errors.content.message}</p>
          )}
        </div>

        {/* Submit Button */}
        <div className="flex justify-end">
          <button
            type="submit"
            disabled={isLoading}
            className="px-6 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isLoading ? "Publishing..." : "Publish Article"}
          </button>
        </div>
      </form>
    </div>
  );
};

export default AddArticle;
