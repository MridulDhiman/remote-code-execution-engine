"use client"

import { z } from "zod"

export const formSchema = z.object({
  code: z.string().min(1),
  stdin: z.string(),
})
