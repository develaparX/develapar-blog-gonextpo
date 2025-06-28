import React, { useState } from "react";

// Interface untuk props TagsInput
interface TagsInputProps {
  tags: string[];
  setTags: React.Dispatch<React.SetStateAction<string[]>>;
  placeholder: string;
  maxTags?: number;
}

export const TagsInput: React.FC<TagsInputProps> = ({
  tags,
  setTags,
  placeholder,
  maxTags = 10,
}) => {
  const [input, setInput] = useState<string>("");

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    setInput(e.target.value);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>): void => {
    if (e.key === "Enter" && input.trim() !== "") {
      e.preventDefault();

      // Check if tag already exists
      if (tags.includes(input.trim())) {
        return;
      }

      // Check if max tags limit reached
      if (tags.length >= maxTags) {
        return;
      }

      setTags([...tags, input.trim()]);
      setInput("");
    }
  };

  const removeTag = (tagToRemove: string): void => {
    setTags(tags.filter((tag: string) => tag !== tagToRemove));
  };

  return (
    <div className="flex flex-wrap items-center p-2 border border-gray-300 rounded-md focus-within:ring-2 focus-within:ring-indigo-500 focus-within:border-indigo-500">
      {tags.map((tag: string, index: number) => (
        <div
          key={index}
          className="flex items-center bg-indigo-100 text-indigo-800 text-sm rounded-md px-2 py-1 m-1"
        >
          <span>{tag}</span>
          <button
            type="button"
            onClick={() => removeTag(tag)}
            className="ml-1 text-indigo-600 hover:text-indigo-800 focus:outline-none"
          >
            &times;
          </button>
        </div>
      ))}
      <input
        type="text"
        value={input}
        onChange={handleInputChange}
        onKeyDown={handleKeyDown}
        placeholder={tags.length === 0 ? placeholder : ""}
        className="flex-grow outline-none px-2 py-1 min-w-[120px]"
      />
    </div>
  );
};
