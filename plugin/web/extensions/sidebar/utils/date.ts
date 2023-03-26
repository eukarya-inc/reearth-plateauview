export const formatDateTime = (d: string, t: string) => {
  const date = d
    ?.split(/-|\//)
    ?.map(s => s.padStart(2, "0"))
    ?.join("-");
  const Time = t
    ?.split(/:/)
    ?.map(s => s.padStart(2, "0"))
    ?.join(":");
  const dateStr = [date, Time].filter(s => !!s).join("T");

  try {
    return new Date(dateStr).toISOString();
  } catch {
    return new Date().toISOString();
  }
};
