import { useCallback, useState } from "react";

import { Template } from "../../types";

export default ({
  backendURL,
  backendProjectName,
  backendAccessToken,
  setLoading,
}: {
  backendURL?: string;
  backendProjectName?: string;
  backendAccessToken?: string;
  setLoading?: React.Dispatch<React.SetStateAction<boolean>>;
}) => {
  const [fieldTemplates, setFieldTemplates] = useState<Template[]>([]);

  const handleTemplateAdd = useCallback(async () => {
    if (!backendURL || !backendProjectName || !backendAccessToken) return;
    const res = await fetch(`${backendURL}/sidebar/${backendProjectName}/templates`, {
      headers: {
        authorization: `Bearer ${backendAccessToken}`,
      },
      method: "POST",
      body: JSON.stringify({ type: "field", name: "新しいテンプレート" }),
    });
    if (res.status !== 200) return;
    const newTemplate = await res.json();
    setFieldTemplates(t => [...t, newTemplate]);
    return newTemplate as Template;
  }, [backendURL, backendProjectName, backendAccessToken]);

  const handleTemplateSave = useCallback(
    async (template: Template) => {
      if (!backendURL || !backendProjectName || !backendAccessToken) return;
      setLoading?.(true);
      const res = await fetch(
        `${backendURL}/sidebar/${backendProjectName}/templates/${template.id}`,
        {
          headers: {
            authorization: `Bearer ${backendAccessToken}`,
          },
          method: "PATCH",
          body: JSON.stringify(template),
        },
      );
      if (res.status !== 200) return;
      const updatedTemplate = await res.json();
      setFieldTemplates(t => {
        return t.map(t2 => {
          if (t2.id === updatedTemplate.id) {
            return updatedTemplate;
          }
          return t2;
        });
      });
      setLoading?.(false);
    },
    [backendURL, backendProjectName, backendAccessToken, setLoading],
  );

  const handleTemplateRemove = useCallback(
    async (id: string) => {
      if (!backendURL || !backendProjectName || !backendAccessToken) return;
      const res = await fetch(`${backendURL}/sidebar/${backendProjectName}/templates/${id}`, {
        headers: {
          authorization: `Bearer ${backendAccessToken}`,
        },
        method: "DELETE",
      });
      if (res.status !== 200) return;
      setFieldTemplates(t => t.filter(t2 => t2.id !== id));
    },
    [backendURL, backendProjectName, backendAccessToken],
  );
  return {
    fieldTemplates,
    setFieldTemplates,
    handleTemplateAdd,
    handleTemplateSave,
    handleTemplateRemove,
  };
};
