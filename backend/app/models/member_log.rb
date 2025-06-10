class MemberLog
  include Mongoid::Document
  include Mongoid::Timestamps

  field :action, type: String  

  field :member, type: Hash   
  field :before_update, type: Hash
  field :after_update, type: Hash
end
