import { NextRequest, NextResponse } from 'next/server'
import { promises as fs } from 'fs'
import path from 'path'
import type { StagesData } from '@/types/stage'

const dataPath = path.join(process.cwd(), 'public/data/stages.json')

export async function GET() {
  try {
    const data = await fs.readFile(dataPath, 'utf8')
    return NextResponse.json(JSON.parse(data))
  } catch (error) {
    console.error('Erreur lecture:', error)
    return NextResponse.json({ stages: [] }, { status: 500 })
  }
}

export async function POST(request: NextRequest) {
  try {
    const body: StagesData = await request.json()
    
    // Validation basique
    if (!body.stages || !Array.isArray(body.stages)) {
      return NextResponse.json({ error: 'Format invalide' }, { status: 400 })
    }

    // Sauvegarde
    await fs.writeFile(dataPath, JSON.stringify(body, null, 2))
    
    return NextResponse.json({ success: true })
  } catch (error) {
    console.error('Erreur sauvegarde:', error)
    return NextResponse.json({ error: 'Erreur serveur' }, { status: 500 })
  }
}
