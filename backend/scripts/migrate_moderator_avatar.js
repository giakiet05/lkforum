// Migration script to convert Moderator.Avatar from string to Image object
// Run this in MongoDB shell or MongoDB Compass

// Use your database
use('LKForum');

// Find all communities with moderators that have string avatars
const communitiesWithOldFormat = db.communities.find({
  'moderators.0.avatar': { $type: 'string' }
}).toArray();

console.log(`Found ${communitiesWithOldFormat.length} communities to migrate`);

let migratedCount = 0;
let errorCount = 0;

communitiesWithOldFormat.forEach(community => {
  try {
    // Convert each moderator's avatar from string to object
    const updatedModerators = community.moderators.map(mod => {
      if (typeof mod.avatar === 'string') {
        return {
          ...mod,
          avatar: mod.avatar ? {
            url: mod.avatar,
            public_id: '',
            uploaded_at: new Date()
          } : null
        };
      }
      return mod; // Already migrated
    });

    // Update the community
    const result = db.communities.updateOne(
      { _id: community._id },
      { $set: { moderators: updatedModerators } }
    );

    if (result.modifiedCount > 0) {
      console.log(`✅ Migrated community: ${community._id} (${community.name})`);
      migratedCount++;
    }
  } catch (error) {
    console.error(`❌ Error migrating community ${community._id}:`, error);
    errorCount++;
  }
});

console.log('\n=== Migration Summary ===');
console.log(`✅ Successfully migrated: ${migratedCount} communities`);
console.log(`❌ Errors: ${errorCount}`);
console.log('Migration completed!');
